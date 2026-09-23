package websocket

import (
	"net/http"
	"sync"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	manager  *connection.Manager
	registry *connection.Registry
	engine   *engine.Engine
	writeMu sync.Mutex
}

func New(
	e *engine.Engine,
	registry *connection.Registry,
) *Server {
	return &Server{
		manager:  connection.NewManager(),
		registry: registry,
		engine:   e,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
