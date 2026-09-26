package bootstrap

import (
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/dispatcher"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/astrahost/astrahost-tunnel/internal/core/handshake"
	"github.com/astrahost/astrahost-tunnel/internal/core/heartbeat"
	"github.com/astrahost/astrahost-tunnel/internal/core/httptunnel"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func NewEngine(
	registry *connection.Registry,
	handshakeConfig handshake.Config,
) *engine.Engine {
	d := dispatcher.New()
	h := handshake.New(registry, handshakeConfig)
	hb := heartbeat.New()
	httpResponse := httptunnel.NewResponseHandler()

	d.Register(protocol.PacketConnect, h.HandlePacket)
	d.Register(protocol.PacketPing, hb.HandlePing)
	d.Register(protocol.PacketPong, hb.HandlePong)
	d.Register(protocol.PacketHTTPResponse, httpResponse.Handle)

	return engine.New(d)
}
