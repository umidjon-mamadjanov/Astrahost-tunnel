package client

import (
	"log"
	"net/url"

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
		return err
	}

	a.Conn = conn
	defer a.Conn.Close()

	log.Println("Connected:", u.String())

	return nil
}
