package connection

import (
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
	ID string

	Conn *websocket.Conn

	State State

	ProtocolVersion uint8
	ClientVersion   string
	Platform        string
	Architecture    string

	CreatedAt time.Time
	LastSeen  time.Time
}
