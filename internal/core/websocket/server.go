package websocket

import (
	"net/http"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/gorilla/websocket"
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
