package client

import (
	"fmt"
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/gorilla/websocket"
)

const (
	initialReconnectDelay = 1 * time.Second
	maxReconnectDelay     = 30 * time.Second
)

type App struct {
	Server string

	LocalHost string
	LocalPort uint16

	TunnelName string

	Conn *websocket.Conn
}

func New(localHost string, localPort uint16) *App {
	return &App{
		Server:     "ws://localhost:7000/connect",
		LocalHost:  localHost,
		LocalPort:  localPort,
		TunnelName: "local",
	}
}

func (a *App) Run() error {
	log.Println("Astra Tunnel Client")

	log.Printf(
		"Local: http://%s:%d",
		a.LocalHost,
		a.LocalPort,
	)

	reconnectDelay := initialReconnectDelay

	for {
		err := a.connect()

		if err == nil {
			reconnectDelay = initialReconnectDelay
			continue
		}

		log.Printf("Connection lost: %v", err)
		log.Printf("Reconnecting in %s...", reconnectDelay)

		time.Sleep(reconnectDelay)

		reconnectDelay *= 2

		if reconnectDelay > maxReconnectDelay {
			reconnectDelay = maxReconnectDelay
		}
	}
}

func (a *App) connect() error {
	u, err := url.Parse(a.Server)
	if err != nil {
		return fmt.Errorf("parse server URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(
		u.String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	a.Conn = conn

	httpHandler := NewHTTPHandler()

	defer func() {
		_ = conn.Close()
		a.Conn = nil
	}()

	log.Println("Connected:", u.String())

	// ------------------------------------------------------------
	// CONNECT
	// ------------------------------------------------------------

	connectPacket, err := protocol.NewConnectPacket(
		protocol.ConnectRequest{
			ProtocolVersion: protocol.Version,
			ClientVersion:   "1.0.0",
			Platform:        runtime.GOOS,
			Architecture:    runtime.GOARCH,
			TunnelName:      a.TunnelName,
			LocalHost:       a.LocalHost,
			LocalPort:       a.LocalPort,
		},
	)
	if err != nil {
		return fmt.Errorf("create CONNECT packet: %w", err)
	}

	connectData, err := protocol.EncodePacket(connectPacket)
	if err != nil {
		return fmt.Errorf("encode CONNECT packet: %w", err)
	}

	if err := conn.WriteMessage(
		websocket.BinaryMessage,
		connectData,
	); err != nil {
		return fmt.Errorf("send CONNECT packet: %w", err)
	}

	log.Printf(
		"CONNECT sent | Local: http://%s:%d",
		a.LocalHost,
		a.LocalPort,
	)

	// ------------------------------------------------------------
	// CONNECT_OK
	// ------------------------------------------------------------

	messageType, data, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read CONNECT_OK: %w", err)
	}

	if messageType != websocket.BinaryMessage {
		return fmt.Errorf(
			"expected binary WebSocket message, got %d",
			messageType,
		)
	}

	responsePacket, err := protocol.DecodePacket(data)
	if err != nil {
		return fmt.Errorf("decode CONNECT_OK packet: %w", err)
	}

	if responsePacket.Header.Type != protocol.PacketConnectOK {
		return fmt.Errorf(
			"expected CONNECT_OK packet, got %d",
			responsePacket.Header.Type,
		)
	}

	response, err := protocol.DecodeConnectResponse(
		responsePacket.Payload,
	)
	if err != nil {
		return fmt.Errorf(
			"decode CONNECT_OK payload: %w",
			err,
		)
	}

	log.Printf(
		"Handshake completed | Session: %s | Tunnel: %s | Heartbeat: %ds",
		response.SessionID,
		response.TunnelID,
		response.HeartbeatInterval,
	)

	if response.PublicURL != "" {
		log.Println()
		log.Println("✓ Tunnel connected")
		log.Printf("Local:  http://%s:%d", a.LocalHost, a.LocalPort)
		log.Printf("Public: %s", response.PublicURL)
		log.Printf("Tunnel: %s", response.TunnelID)
		log.Println()
	}

	// ------------------------------------------------------------
	// PACKET LOOP
	// ------------------------------------------------------------

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf(
				"connection closed: %w",
				err,
			)
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		packet, err := protocol.DecodePacket(data)
		if err != nil {
			return fmt.Errorf(
				"decode packet: %w",
				err,
			)
		}

		switch packet.Header.Type {

		// --------------------------------------------------------
		// PING
		// --------------------------------------------------------

		case protocol.PacketPing:
			if err := protocol.ValidatePing(packet); err != nil {
				return fmt.Errorf(
					"invalid PING: %w",
					err,
				)
			}

			pong := protocol.NewPongPacket(
				packet.Header.RequestID,
			)

			pongData, err := protocol.EncodePacket(pong)
			if err != nil {
				return fmt.Errorf(
					"encode PONG: %w",
					err,
				)
			}

			if err := conn.WriteMessage(
				websocket.BinaryMessage,
				pongData,
			); err != nil {
				return fmt.Errorf(
					"send PONG: %w",
					err,
				)
			}

			log.Println("PING received → PONG sent")

		// --------------------------------------------------------
		// HTTP REQUEST
		// --------------------------------------------------------

		case protocol.PacketHTTPRequest:
			responsePacket, err := httpHandler.Handle(packet)
			if err != nil {
				return fmt.Errorf(
					"handle HTTP_REQUEST: %w",
					err,
				)
			}

			responseData, err := protocol.EncodePacket(
				responsePacket,
			)
			if err != nil {
				return fmt.Errorf(
					"encode HTTP_RESPONSE: %w",
					err,
				)
			}

			if err := conn.WriteMessage(
				websocket.BinaryMessage,
				responseData,
			); err != nil {
				return fmt.Errorf(
					"send HTTP_RESPONSE: %w",
					err,
				)
			}

			log.Printf(
				"HTTP_REQUEST forwarded → HTTP_RESPONSE sent | Request: %s",
				packet.Header.RequestID,
			)

		// --------------------------------------------------------
		// UNKNOWN / OTHER PACKETS
		// --------------------------------------------------------

		default:
			log.Printf(
				"Received packet: type=%d",
				packet.Header.Type,
			)
		}
	}
}
