package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/engine"
	"github.com/astrahost/astrahost-tunnel/internal/core/heartbeat"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader       websocket.Upgrader
	manager        *connection.Manager
	registry       *connection.Registry
	engine         *engine.Engine
	writeMu        sync.Mutex
	timeoutChecker *heartbeat.TimeoutChecker
}

func New(
	e *engine.Engine,
	registry *connection.Registry,
) *Server {
	manager := connection.NewManager()

	server := &Server{
		manager:  manager,
		registry: registry,
		engine:   e,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	server.timeoutChecker = heartbeat.NewTimeoutChecker(
		manager,
		registry,
	)

	go server.runTimeoutChecker()

	return server
}

func (s *Server) runTimeoutChecker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for now := range ticker.C {
		s.timeoutChecker.Check(now)
	}
}

// SendPacket sends a protocol packet to a connected client session.
func (s *Server) SendPacket(
	session *connection.Session,
	packet protocol.Packet,
) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if session.Conn == nil {
		return fmt.Errorf("session connection is nil")
	}

	data, err := protocol.EncodePacket(packet)
	if err != nil {
		return fmt.Errorf("encode packet: %w", err)
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	if err := session.Conn.WriteMessage(
		websocket.BinaryMessage,
		data,
	); err != nil {
		return fmt.Errorf("write WebSocket packet: %w", err)
	}

	return nil
}
