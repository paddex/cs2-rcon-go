package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"paddex.net/cs2-rcon-go/errors"
)

type Client struct {
	conn net.Conn
	id   int32
}

type pkg struct {
	size    int32
	id      int32
	pkgType int32
	body    []byte
}

func NewClient(addr string, port int) (*Client, error) {
	serverAddr := fmt.Sprintf("%s:%d", addr, port)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, id: 1}, nil
}

func (c *Client) Auth(password string) error {
	pw := []byte(password)
	id := c.id
	c.id = c.id + 1
	authPkg := pkg{
		size:    int32(len(pw)) + 10,
		id:      id,
		pkgType: 3,
		body:    append(pw, 0x00, 0x00),
	}

	err := c.send(authPkg)
	if err != nil {
		return fmt.Errorf("Error while sending AuthPkg: %e", err)
	}

	replyPkg, err := c.receive()
	if err != nil {
		return fmt.Errorf("Error while reading AuthReply: %e", err)
	}

	if replyPkg.id != id {
		return errors.AuthError{Msg: "Bad rcon password"}
	}

	return nil
}

func (c *Client) Exec(cmd string) (string, error) {
	pw := []byte(cmd)
	id := c.id
	c.id = c.id + 1
	cmdPkg := pkg{
		size:    int32(len(pw)) + 10,
		id:      id,
		pkgType: 2,
		body:    append(pw, 0x00, 0x00),
	}

	err := c.send(cmdPkg)
	if err != nil {
		return "", fmt.Errorf("Error while sending CmdPkg: %e", err)
	}

	replyPkg, err := c.receive()
	if err != nil {
		return "", fmt.Errorf("Error while reading CmdReply: %e", err)
	}

	reply := string(replyPkg.body)

	return reply, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) send(pkg pkg) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, pkg.size)
	binary.Write(buf, binary.LittleEndian, pkg.id)
	binary.Write(buf, binary.LittleEndian, pkg.pkgType)

	buf.Write(pkg.body)

	_, err := buf.WriteTo(c.conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) receive() (*pkg, error) {
	var pkg pkg

	buf := make([]byte, 4)
	readTotal := 0

	// read size
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	readTotal += n
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &pkg.size)
	if err != nil {
		return nil, err
	}

	// read id
	n, err = c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	readTotal += n
	reader = bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &pkg.id)
	if err != nil {
		return nil, err
	}

	// read type
	n, err = c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	readTotal += n
	reader = bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &pkg.pkgType)
	if err != nil {
		return nil, err
	}

	// read body
	if pkg.size > 4096 {
		for readTotal < int(pkg.size) {
			bodyLen := 4096
			buf = make([]byte, bodyLen)
			n, err := c.conn.Read(buf)
			if err != nil {
				return nil, err
			}
			readTotal += n
			pkg.body = append(pkg.body, buf...)
		}
	} else {
		bodyLen := pkg.size - 8
		buf = make([]byte, bodyLen)
		_, err = c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		pkg.body = buf
	}

	return &pkg, nil
}
