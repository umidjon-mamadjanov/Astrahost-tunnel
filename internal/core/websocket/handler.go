package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
	"github.com/astrahost/astrahost-tunnel/internal/core/transport"
)

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	id := uuid.NewString()

	session := &connection.Session{
		ID:        id,
		Conn:      conn,
		Transport: transport.NewWebSocket(conn),

		State: connection.StateConnected,

		Send: make(chan *protocol.Packet, 32),

		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	s.manager.Add(session)

	log.Printf("Client connected: %s | Online: %d",
		id,
		s.manager.Count(),
	)

	// ===== Clientdan xabar kutish =====

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// ===== Sessionni o'chirish =====

	s.manager.Remove(id)

	log.Printf("Client disconnected: %s | Online: %d",
		id,
		s.manager.Count(),
	)
}
