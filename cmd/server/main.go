package main

import (
	"log"

	"github.com/astrahost/astrahost-tunnel/internal/server"
)

func main() {

	app := server.New()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
