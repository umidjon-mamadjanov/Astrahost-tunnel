package transport

import "github.com/astrahost/astrahost-tunnel/internal/core/protocol"

type Transport interface {
	Read() (*protocol.Packet, error)

	Write(*protocol.Packet) error

	Close() error
}
