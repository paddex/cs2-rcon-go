package gameinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Client struct {
	conn    net.Conn
	lastPkg pkg
}

type pkg struct {
	header    int32
	pkgtype   byte
	payload   []byte
	challenge []byte
}

func NewClient(addr string, port int) (*Client, error) {
	serverAddr := fmt.Sprintf("%s:%d", addr, port)

	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) send(pkg pkg) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, pkg.header)
	binary.Write(buf, binary.LittleEndian, pkg.pkgtype)
	if len(pkg.payload) != 0 {
		buf.Write(pkg.payload)
	}
	if len(pkg.challenge) != 0 {
		buf.Write(pkg.challenge)
	}

	_, err := buf.WriteTo(c.conn)
	if err != nil {
		return err
	}

	c.lastPkg = pkg

	return nil
}

func (c *Client) read() ([]byte, error) {
	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	// Check if challenged && repeat request with challenge number
	if buf[4] == 0x41 {
		pkg := c.lastPkg
		pkg.challenge = buf[5:n]

		err = c.send(pkg)
		if err != nil {
			return nil, err
		}

		n, err = c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
	}

	return buf[4:n], nil
}

func (c *Client) Close() {
	c.conn.Close()
}
