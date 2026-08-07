package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func Encode(packet Packet) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, packet.Header.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint8(packet.Header.Type)); err != nil {
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

	if len(data) < 6 {
		return Packet{}, errors.New("packet too small")
	}

	packet := Packet{}

	packet.Header.Version = data[0]
	packet.Header.Type = PacketType(data[1])
	packet.Header.Length = binary.BigEndian.Uint32(data[2:6])

	if len(data) != int(packet.Header.Length)+6 {
		return Packet{}, errors.New("invalid payload length")
	}

	packet.Payload = data[6:]

	return packet, nil
}
