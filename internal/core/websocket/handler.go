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
		ID:              id,
		Conn:            conn,
		State:           connection.StateConnected,
		CreatedAt:       time.Now(),
		LastSeen:        time.Now(),
		PendingRequests: connection.NewPendingRequests(),
	}

	s.manager.Add(session)

	log.Printf(
		"Client connected: %s | Online: %d",
		id,
		s.manager.Count(),
	)

	defer func() {
		tunnelID := session.GetTunnelID()

		if tunnelID != "" {
			s.registry.Unregister(tunnelID)

			log.Printf(
				"Tunnel unregistered: %s | Session: %s",
				tunnelID,
				id,
			)
		}

		session.SetState(connection.StateClosed)

		s.manager.Remove(id)

		log.Printf(
			"Client disconnected: %s | Online: %d",
			id,
			s.manager.Count(),
		)
	}()

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			ping := protocol.NewPingPacket()

			data, err := protocol.EncodePacket(ping)
			if err != nil {
				log.Printf(
					"Failed to encode PING | Session: %s | Error: %v",
					id,
					err,
				)
				return
			}

			s.writeMu.Lock()
			err = conn.WriteMessage(websocket.BinaryMessage, data)
			s.writeMu.Unlock()

			if err != nil {
				log.Printf(
					"Failed to send PING | Session: %s | Error: %v",
					id,
					err,
				)
				return
			}

			log.Printf("PING sent | Session: %s", id)
		}
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

		session.Touch(time.Now())

		response, err := s.engine.Handle(session, &packet)
		if err != nil {
			log.Printf(
				"Packet handling failed | Session: %s | Type: %d | Error: %v",
				id,
				packet.Header.Type,
				err,
			)
			continue
		}

		if response != nil {
			data, err := protocol.EncodePacket(*response)
			if err != nil {
				log.Printf(
					"Failed to encode response | Session: %s | Error: %v",
					id,
					err,
				)
				continue
			}

			s.writeMu.Lock()
			err = conn.WriteMessage(websocket.BinaryMessage, data)
			s.writeMu.Unlock()

			if err != nil {
				log.Printf(
					"Failed to send response | Session: %s | Error: %v",
					id,
					err,
				)
				return
			}
		}
	}
}
