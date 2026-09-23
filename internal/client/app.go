package client

import (
	"fmt"
	"log"
	"net/url"
	"runtime"

	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/gorilla/websocket"
)

type App struct {
	Server string
	Conn   *websocket.Conn
}

func New() *App {
	return &App{
		Server: "ws://localhost:7000/connect",
	}
}

func (a *App) Run() error {
	log.Println("Astra Tunnel Client")

	u, err := url.Parse(a.Server)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(
		u.String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	a.Conn = conn
	defer a.Conn.Close()

	log.Println("Connected:", u.String())

	connectPacket, err := protocol.NewConnectPacket(
		protocol.ConnectRequest{
			ProtocolVersion: protocol.Version,
			ClientVersion:   "1.0.0",
			Platform:        runtime.GOOS,
			Architecture:    runtime.GOARCH,
			TunnelName:      "local",
		},
	)
	if err != nil {
		return fmt.Errorf("create CONNECT packet: %w", err)
	}

	connectData, err := protocol.EncodePacket(connectPacket)
	if err != nil {
		return fmt.Errorf("encode CONNECT packet: %w", err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, connectData); err != nil {
		return fmt.Errorf("send CONNECT packet: %w", err)
	}

	log.Println("CONNECT sent")

	messageType, data, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read CONNECT_OK: %w", err)
	}

	if messageType != websocket.BinaryMessage {
		return fmt.Errorf("expected binary WebSocket message, got %d", messageType)
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

	response, err := protocol.DecodeConnectResponse(responsePacket.Payload)
	if err != nil {
		return fmt.Errorf("decode CONNECT_OK payload: %w", err)
	}

	log.Printf(
		"Handshake completed | Session: %s | Tunnel: %s | Heartbeat: %ds",
		response.SessionID,
		response.TunnelID,
		response.HeartbeatInterval,
	)

	log.Println("Handshake completed")

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("connection closed: %w", err)
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		packet, err := protocol.DecodePacket(data)
		if err != nil {
			return fmt.Errorf("decode packet: %w", err)
		}

		switch packet.Header.Type {
		case protocol.PacketPing:
			if err := protocol.ValidatePing(packet); err != nil {
				return fmt.Errorf("invalid PING: %w", err)
			}

			pong := protocol.NewPongPacket(packet.Header.RequestID)

			pongData, err := protocol.EncodePacket(pong)
			if err != nil {
				return fmt.Errorf("encode PONG: %w", err)
			}

			if err := conn.WriteMessage(websocket.BinaryMessage, pongData); err != nil {
				return fmt.Errorf("send PONG: %w", err)
			}

			log.Println("PING received → PONG sent")

		default:
			log.Printf("Received packet: type=%d", packet.Header.Type)
		}
	}
}
