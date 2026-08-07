package client

import "log"

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {

	log.Println("Astra Tunnel Client")

	return nil
}
