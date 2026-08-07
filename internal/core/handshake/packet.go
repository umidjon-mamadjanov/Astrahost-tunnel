package handshake

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
)

func (h *Handler) HandlePacket(
	session *connection.Session,
	packet *protocol.Packet,
) error {

	var req protocol.ConnectRequest

	if err := protocol.DecodeJSON(packet.Payload, &req); err != nil {
		return err
	}

	if req.ProtocolVersion != protocol.ProtocolVersion {
		return fmt.Errorf(
			"unsupported protocol version: client=%d server=%d",
			req.ProtocolVersion,
			protocol.ProtocolVersion,
		)
	}

	session.ProtocolVersion = req.ProtocolVersion
	session.ClientVersion = req.ClientVersion
	session.Platform = req.Platform
	session.Architecture = req.Architecture
	session.State = connection.StateReady

	return nil
}
