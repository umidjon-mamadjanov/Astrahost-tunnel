package main

import (
	"log"
	"net/http"

	"github.com/astrahost/server/internal/tunnel"
	ws "github.com/astrahost/server/internal/websocket"
)

func main() {

	registry := tunnel.NewRegistry()

	http.HandleFunc("/connect", ws.NewHandler(registry))

	log.Println("Astra Tunnel Server started")
	log.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
