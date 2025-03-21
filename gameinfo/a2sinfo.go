package gameinfo

type A2S_Info_Response struct {
	Header     string
	Protocol   byte
	Name       string
	MapName    string
	Folder     string
	Game       string
	Id         int16
	Players    int8
	MaxPlayers int8
	Bots       int8
	ServerType string
	Env        string
	Visibility int8
	Vac        int8
	Version    string
	Edf        byte
	Port       int16
	Steamid    uint64
	SpecPort   int16
	SpecName   string
	Keywords   string
	Gameid     uint64
}

func (c *Client) GetServerInfo() (*A2S_Info_Response, error) {
	info := pkg{
		header:    -1,
		pkgtype:   'T',
		payload:   append([]byte("Source Engine Query"), 0x0),
		challenge: []byte{},
	}

	err := c.send(info)
	if err != nil {
		return nil, err
	}

	buf, err := c.read()
	if err != nil {
		return nil, err
	}

	response := A2S_Info_Response{}
	response.Header, buf = readChar(buf)
	response.Protocol, buf = readByte(buf)
	response.Name, buf = readString(buf)
	response.MapName, buf = readString(buf)
	response.Folder, buf = readString(buf)
	response.Game, buf = readString(buf)
	response.Id, buf, err = readInt16(buf)
	if err != nil {
		return nil, err
	}
	response.Players, buf, err = readInt8(buf)
	if err != nil {
		return nil, err
	}
	response.MaxPlayers, buf, err = readInt8(buf)
	if err != nil {
		return nil, err
	}
	response.Bots, buf, err = readInt8(buf)
	if err != nil {
		return nil, err
	}
	response.ServerType, buf = readChar(buf)
	response.Env, buf = readChar(buf)
	response.Visibility, buf, err = readInt8(buf)
	if err != nil {
		return nil, err
	}
	response.Vac, buf, err = readInt8(buf)
	if err != nil {
		return nil, err
	}
	response.Version, buf = readString(buf)
	response.Edf, buf = readByte(buf)

	if (response.Edf & 0x80) != 0x0 {
		response.Port, buf, err = readInt16(buf)
	}
	if (response.Edf & 0x10) != 0x0 {
		response.Steamid, buf, err = readUint64(buf)
	}
	if (response.Edf & 0x40) != 0x0 {
		response.SpecPort, buf, err = readInt16(buf)
		if err != nil {
			return nil, err
		}
		response.SpecName, buf = readString(buf)
	}
	if (response.Edf & 0x20) != 0x0 {
		response.Keywords, buf = readString(buf)
	}
	if (response.Edf & 0x01) != 0x0 {
		response.Gameid, buf, err = readUint64(buf)
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}
