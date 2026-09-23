package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	id := uuid.NewString()

	session := &connection.Session{
		ID:        id,
		Conn:      conn,
		State:     connection.StateConnected,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	s.manager.Add(session)

	log.Printf(
		"Client connected: %s | Online: %d",
		id,
		s.manager.Count(),
	)

	defer func() {
		s.manager.Remove(id)

		log.Printf(
			"Client disconnected: %s | Online: %d",
			id,
			s.manager.Count(),
		)
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType != websocket.BinaryMessage {
			log.Printf(
				"Ignoring non-binary WebSocket message from %s",
				id,
			)
			continue
		}

		packet, err := protocol.DecodePacket(data)
		if err != nil {
			log.Printf(
				"Invalid packet from %s: %v",
				id,
				err,
			)
			continue
		}

		session.LastSeen = time.Now()

		if err := s.engine.Handle(session, &packet); err != nil {
			log.Printf(
				"Packet handling failed | Session: %s | Type: %d | Error: %v",
				id,
				packet.Header.Type,
				err,
			)
		}
	}
}
