package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHandler(registry *tunnel.Registry) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// registry shu yerda ishlatiladi

	}
}
