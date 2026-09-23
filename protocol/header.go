package protocol

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
)

type Header struct {
	Version       uint8
	Type          PacketType
	Flags         uint16
	RequestID     uuid.UUID
	PayloadLength uint32
}

func (h Header) Encode() ([]byte, error) {
	buf := make([]byte, HeaderSize)

	buf[0] = h.Version
	buf[1] = byte(h.Type)

	binary.BigEndian.PutUint16(buf[2:4], h.Flags)

	copy(buf[4:20], h.RequestID[:])

	binary.BigEndian.PutUint32(buf[20:24], h.PayloadLength)

	return buf, nil
}

func DecodeHeader(data []byte) (Header, error) {
	if len(data) < HeaderSize {
		return Header{}, fmt.Errorf("invalid header length: got %d, want %d", len(data), HeaderSize)
	}

	var requestID uuid.UUID
	copy(requestID[:], data[4:20])

	return Header{
		Version:       data[0],
		Type:          PacketType(data[1]),
		Flags:         binary.BigEndian.Uint16(data[2:4]),
		RequestID:     requestID,
		PayloadLength: binary.BigEndian.Uint32(data[20:24]),
	}, nil
}

func (h Header) Validate() error {
	if h.Version != Version {
		return fmt.Errorf("unsupported protocol version: %d", h.Version)
	}

	if h.PayloadLength > MaxPayloadLength {
		return fmt.Errorf(
			"payload too large: %d bytes, maximum %d",
			h.PayloadLength,
			MaxPayloadLength,
		)
	}

	return nil
}
