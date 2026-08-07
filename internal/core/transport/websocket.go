package transport

import (
	"github.com/gorilla/websocket"

	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
)

type WebSocketTransport struct {
	conn *websocket.Conn
}

func NewWebSocket(
	conn *websocket.Conn,
) *WebSocketTransport {

	return &WebSocketTransport{
		conn: conn,
	}
}

func (t *WebSocketTransport) Read() (*protocol.Packet, error) {
	panic("not implemented")
}

func (t *WebSocketTransport) Write(
	packet *protocol.Packet,
) error {
	panic("not implemented")
}

func (t *WebSocketTransport) Close() error {
	return t.conn.Close()
}
