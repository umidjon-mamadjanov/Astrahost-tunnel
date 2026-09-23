package handshake

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func (h *Handler) HandlePacket(
	session *connection.Session,
	packet *protocol.Packet,
) (*protocol.Packet, error) {
	if packet.Header.Type != protocol.PacketConnect {
		return nil, fmt.Errorf(
			"expected CONNECT packet, got %d",
			packet.Header.Type,
		)
	}

	req, err := protocol.DecodeConnectRequest(packet.Payload)
	if err != nil {
		return nil, fmt.Errorf("decode CONNECT: %w", err)
	}

	resp, err := h.HandleConnect(session, req)
	if err != nil {
		return nil, err
	}

	responsePacket, err := protocol.NewConnectOKPacket(*resp)
	if err != nil {
		return nil, fmt.Errorf("create CONNECT_OK: %w", err)
	}

	responsePacket.Header.RequestID = packet.Header.RequestID

	return &responsePacket, nil
}
