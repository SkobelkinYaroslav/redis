package server

import (
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"net"
	"redis/command"
)

type Peer struct {
	conn    net.Conn
	msgChan chan Message

	delChan chan *Peer
}

func NewPeer(conn net.Conn, msgCh chan Message, delChan chan *Peer) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgCh,
		delChan: delChan,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) Read() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delChan <- p
			break
		}
		if err != nil {
			return err
		}
		if v.Type() == resp.Array {
			var cmd command.Command
			for _, val := range v.Array() {
				switch val.String() {
				case command.Set:
					if len(v.Array()) != 3 {
						return fmt.Errorf("command Parse: invalid number of variables in command")
					}
					cmd = command.SetCommand{
						Key: v.Array()[1].Bytes(),
						Val: v.Array()[2].Bytes(),
					}

				case command.Get:
					if len(v.Array()) != 2 {
						return fmt.Errorf("command Parse: invalid number of variables in command")
					}
					cmd = command.GetCommand{
						Key: v.Array()[1].Bytes(),
					}
				}
			}
			p.msgChan <- Message{
				cmd:    cmd,
				client: p,
			}
		}
	}
	return fmt.Errorf("command Parse: cant parse string")
}
