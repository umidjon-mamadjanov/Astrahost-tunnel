package bootstrap

import (
	"github.com/astrahost/astrahost-tunnel/internal/core/dispatcher"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/astrahost/astrahost-tunnel/internal/core/handshake"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func NewEngine() *engine.Engine {
	d := dispatcher.New()
	h := handshake.New()

	d.Register(protocol.PacketConnect, h.HandlePacket)

	return engine.New(d)
}
