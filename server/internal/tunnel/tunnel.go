package tunnel

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Tunnel struct {
	ID         uuid.UUID
	Subdomain  string
	Conn       *websocket.Conn
	Connected  time.Time
	LastSeen   time.Time
}
