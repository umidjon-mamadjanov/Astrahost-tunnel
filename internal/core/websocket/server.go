package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

type Server struct {
	upgrader websocket.Upgrader
	manager  *connection.Manager
}

func New() *Server {
	return &Server{
		manager: connection.NewManager(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
