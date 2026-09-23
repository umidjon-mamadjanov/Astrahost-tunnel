package server

import (
	"log"
	"net/http"

	"github.com/astrahost/astrahost-tunnel/internal/core/bootstrap"
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/websocket"
)

type App struct {
	cfg Config
	ws  *websocket.Server
}

func New() *App {
	registry := connection.NewRegistry()
	engine := bootstrap.NewEngine(registry)

	return &App{
		cfg: DefaultConfig(),
		ws:  websocket.New(engine, registry),
	}
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:    a.cfg.Address,
		Handler: a.routes(),
	}

	log.Printf("Astra Tunnel Server listening on %s", a.cfg.Address)

	return server.ListenAndServe()
}
