package handshake

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func (h *Handler) HandlePacket(
	session *connection.Session,
	packet *protocol.Packet,
) error {
	if packet.Header.Type != protocol.PacketConnect {
		return fmt.Errorf(
			"expected CONNECT packet, got %d",
			packet.Header.Type,
		)
	}

	req, err := protocol.DecodeConnectRequest(packet.Payload)
	if err != nil {
		return fmt.Errorf("decode CONNECT: %w", err)
	}

	_, err = h.HandleConnect(session, req)
	if err != nil {
		return err
	}

	return nil
}
