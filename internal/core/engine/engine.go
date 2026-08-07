package engine

import (
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/dispatcher"
	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
)

type Engine struct {
	dispatcher *dispatcher.Dispatcher
}

func New(d *dispatcher.Dispatcher) *Engine {
	return &Engine{
		dispatcher: d,
	}
}

func (e *Engine) Handle(
	session *connection.Session,
	packet protocol.Packet,
) error {
	return e.dispatcher.Dispatch(session, packet)
}
