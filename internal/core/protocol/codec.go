package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func EncodeJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func DecodeJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func Encode(packet Packet) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, packet.Header.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, packet.Header.Type); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, packet.Header.Length); err != nil {
		return nil, err
	}

	if _, err := buf.Write(packet.Payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decode(data []byte) (Packet, error) {
	reader := bytes.NewReader(data)

	var header Header

	if err := binary.Read(reader, binary.BigEndian, &header.Version); err != nil {
		return Packet{}, err
	}

	if err := binary.Read(reader, binary.BigEndian, &header.Type); err != nil {
		return Packet{}, err
	}

	if err := binary.Read(reader, binary.BigEndian, &header.Length); err != nil {
		return Packet{}, err
	}

	payload := make([]byte, header.Length)

	n, err := reader.Read(payload)
	if err != nil {
		return Packet{}, err
	}

	if uint32(n) != header.Length {
		return Packet{}, fmt.Errorf("invalid payload length")
	}

	return Packet{
		Header:  header,
		Payload: payload,
	}, nil
}
