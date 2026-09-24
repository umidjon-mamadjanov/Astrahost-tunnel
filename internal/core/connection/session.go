package connection

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type State uint8

const (
	StateConnected State = iota
	StateHandshake
	StateReady
	StateClosed
)

type Session struct {
	mu sync.RWMutex

	ID string

	Conn *websocket.Conn

	State State

	ProtocolVersion uint8
	TunnelID        string
	ClientVersion   string
	Platform        string
	Architecture    string

	LocalHost string
	LocalPort uint16

	CreatedAt time.Time
	LastSeen  time.Time

	PendingRequests *PendingRequests
}

func (s *Session) Touch(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.LastSeen = now
}

func (s *Session) GetLastSeen() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.LastSeen
}

func (s *Session) SetState(state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.State = state
}

func (s *Session) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.State
}

func (s *Session) SetTunnelID(tunnelID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TunnelID = tunnelID
}

func (s *Session) GetTunnelID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.TunnelID
}

func (s *Session) SetLocalTarget(host string, port uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.LocalHost = host
	s.LocalPort = port
}

func (s *Session) GetLocalTarget() (string, uint16) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.LocalHost, s.LocalPort
}
