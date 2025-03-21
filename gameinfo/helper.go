package gameinfo

import (
	"bytes"
	"encoding/binary"
	"slices"
)

func readInt8(in []byte) (int8, []byte, error) {
	var res int8
	err := binary.Read(bytes.NewReader(in[:1]), binary.LittleEndian, &res)
	if err != nil {
		return 0, nil, err
	}

	return res, in[1:], nil
}

func readInt16(in []byte) (int16, []byte, error) {
	var res int16
	err := binary.Read(bytes.NewReader(in[:2]), binary.LittleEndian, &res)
	if err != nil {
		return 0, nil, err
	}

	return res, in[2:], nil
}

func readUint64(in []byte) (uint64, []byte, error) {
	var res uint64
	err := binary.Read(bytes.NewReader(in[:8]), binary.LittleEndian, &res)
	if err != nil {
		return 0, nil, err
	}

	return res, in[8:], nil
}

func readChar(in []byte) (string, []byte) {
	char := in[0]

	return string(char), in[1:]
}

func readByte(in []byte) (byte, []byte) {
	return in[0], in[1:]
}

func readString(in []byte) (string, []byte) {
	return string(in[:slices.Index(in, 0x0)]), in[slices.Index(in, 0x0)+1:]
}
