package websocket

import (
	"net/http"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	manager  *connection.Manager
	engine   *engine.Engine
}

func New(e *engine.Engine) *Server {
	return &Server{
		manager: connection.NewManager(),
		engine:  e,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
