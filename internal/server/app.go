package server

import (
	"log"
	"net/http"

	"github.com/astrahost/astrahost-tunnel/internal/core/bootstrap"
	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/handshake"
	"github.com/astrahost/astrahost-tunnel/internal/core/websocket"
)

type App struct {
	cfg      Config
	ws       *websocket.Server
	registry *connection.Registry
}

func New() *App {
	cfg := DefaultConfig()
	registry := connection.NewRegistry()

	engine := bootstrap.NewEngine(
		registry,
		handshake.Config{
			BaseDomain: cfg.BaseDomain,
			Scheme:     cfg.Scheme,
		},
	)

	return &App{
		cfg:      cfg,
		ws:       websocket.New(engine, registry),
		registry: registry,
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
