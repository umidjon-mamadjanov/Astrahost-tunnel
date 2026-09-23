package bootstrap

import (
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/dispatcher"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/astrahost/astrahost-tunnel/internal/core/handshake"
	"github.com/astrahost/astrahost-tunnel/internal/core/heartbeat"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func NewEngine(registry *connection.Registry) *engine.Engine {
	d := dispatcher.New()
	h := handshake.New(registry)
	hb := heartbeat.New()

	d.Register(protocol.PacketConnect, h.HandlePacket)
	d.Register(protocol.PacketPing, hb.HandlePing)
	d.Register(protocol.PacketPong, hb.HandlePong)

	return engine.New(d)
}
