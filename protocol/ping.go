package protocol

import (
	"fmt"

	"github.com/google/uuid"
)

func NewPingPacket() Packet {
	return Packet{
		Header: Header{
			Version:   Version,
			Type:      PacketPing,
			RequestID: uuid.New(),
		},
	}
}

func NewPongPacket(requestID uuid.UUID) Packet {
	return Packet{
		Header: Header{
			Version:   Version,
			Type:      PacketPong,
			RequestID: requestID,
		},
	}
}

func ValidatePing(packet Packet) error {
	if packet.Header.Type != PacketPing {
		return fmt.Errorf("expected PING packet, got %d", packet.Header.Type)
	}

	if len(packet.Payload) != 0 {
		return fmt.Errorf("PING payload must be empty")
	}

	return packet.Header.Validate()
}

func ValidatePong(packet Packet) error {
	if packet.Header.Type != PacketPong {
		return fmt.Errorf("expected PONG packet, got %d", packet.Header.Type)
	}

	if len(packet.Payload) != 0 {
		return fmt.Errorf("PONG payload must be empty")
	}

	return packet.Header.Validate()
}
