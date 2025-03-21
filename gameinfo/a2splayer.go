package gameinfo

import (
	"fmt"
)

type A2S_Player_Response struct{}

func (c *Client) GetPlayers() (*A2S_Player_Response, error) {
	players := pkg{
		header:    -1,
		pkgtype:   'U',
		payload:   []byte{},
		challenge: []byte{0xFF, 0xFF, 0xFF, 0xFF},
	}

	err := c.send(players)
	if err != nil {
		return nil, err
	}

	buf, err := c.read()
	if err != nil {
		return nil, err
	}

	fmt.Printf("HEX: % x\n", buf)
	return nil, nil
}
