package client

import (
	"bytes"
	"context"
	"github.com/tidwall/resp"
	"net"
)

type Client struct {
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value any) error {

	var buf bytes.Buffer

	wr := resp.NewWriter(&buf)

	switch value.(type) {
	case int:
		wr.WriteArray(
			[]resp.Value{
				resp.StringValue("set"),
				resp.StringValue(key),
				resp.IntegerValue(value.(int)),
			},
		)

	case string:
		wr.WriteArray(
			[]resp.Value{
				resp.StringValue("set"),
				resp.StringValue(key),
				resp.StringValue(value.(string)),
			},
		)
	}

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("get"),
			resp.StringValue(key),
		},
	)

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)

	n, err := c.conn.Read(b)

	return string(b[:n]), nil

}

func (c *Client) Close() error {
	return c.conn.Close()
}
