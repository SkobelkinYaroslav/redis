package server

import (
	"fmt"
	"log"
	"net"
	"redis/command"
)

type Config struct {
	ListenAdd string
}

type Message struct {
	cmd    command.Command
	client *Peer
}

type Server struct {
	Config
	ln net.Listener

	Clients   map[*Peer]bool
	addPeerCh chan *Peer
	delPeerCh chan *Peer

	quitCh chan struct{}
	msgCh  chan Message

	kv *command.KV
}

func New(cfg Config) *Server {
	if len(cfg.ListenAdd) == 0 {
		cfg.ListenAdd = "localhost:8080"
	}

	return &Server{
		Config:    cfg,
		Clients:   make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAdd)
	if err != nil {
		return err
	}

	s.ln = ln
	s.kv = command.NewKV()

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) Stop() {
	log.Println("Stop() hit")
	s.quitCh <- struct{}{}
}

func (s *Server) loop() {
	for {
		select {
		case Client := <-s.addPeerCh:
			s.Clients[Client] = true
		case Client := <-s.delPeerCh:
			delete(s.Clients, Client)
		case <-s.quitCh:
			s.ln.Close()
			log.Println("got value from quitChannel")
			return

		case msg := <-s.msgCh:
			err := s.handleRawMessage(msg)
			if err != nil {
				delete(s.Clients, msg.client)
			}
		}

	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
				log.Println("Server stopped accepting new connections")
				return nil
			}
			continue
		}

		s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	client := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- client

	go client.Read()

}

func (s *Server) handleRawMessage(msg Message) error {

	switch v := msg.cmd.(type) {
	case command.SetCommand:
		return s.kv.Set(string(v.Key), v.Val)

	case command.GetCommand:
		val, ok := s.kv.Get(string(v.Key))
		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.client.Send(val)
		if err != nil {
			return err
		}
	}

	return nil
}
