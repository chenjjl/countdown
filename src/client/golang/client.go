package golang

import (
	_const "countdown/src/const"
	"countdown/src/event"
	"encoding/binary"
	"errors"
	uuid "github.com/satori/go.uuid"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

// AddEventOneway add event in oneway, do not wait to receive callback
func (c *Client) AddEventOneway(topic string, body []byte, expiration time.Duration) error {
	_event, err := event.NewEvent(topic, body, uuid.NewV4().String(), expiration)
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(_event.Encode() + "/n"))
	if err != nil {
		return err
	}
	return nil
}

// AddEvent add event in sync
func (c *Client) AddEvent(topic string, body []byte, expiration time.Duration) error {
	_event, err := event.NewEvent(topic, body, uuid.NewV4().String(), expiration)
	if err != nil {
		return err
	}
	_, err = c.conn.Write([]byte(_event.Encode() + "/n"))
	if err != nil {
		return err
	}
	var buf []byte
	_, err = c.conn.Read(buf[:])
	if err != nil {
		return err
	}
	status := binary.BigEndian.Uint32(buf)
	if status == _const.FAIL {
		return errors.New("failed to add time wheel server")
	}
	return nil
}
