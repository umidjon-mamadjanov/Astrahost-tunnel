package dispatcher

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

type HandlerFunc func(
	session *connection.Session,
	packet *protocol.Packet,
) error

type Dispatcher struct {
	handlers map[protocol.PacketType]HandlerFunc
}

func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[protocol.PacketType]HandlerFunc),
	}
}

func (d *Dispatcher) Register(
	packetType protocol.PacketType,
	handler HandlerFunc,
) {
	d.handlers[packetType] = handler
}

func (d *Dispatcher) Dispatch(
	session *connection.Session,
	packet *protocol.Packet,
) error {
	handler, ok := d.handlers[packet.Header.Type]
	if !ok {
		return fmt.Errorf("unknown packet type: %d", packet.Header.Type)
	}

	return handler(session, packet)
}
