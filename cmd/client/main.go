package main

import (
	"log"

	"github.com/astrahost/astrahost-tunnel/internal/client"
)

func main() {

	app := client.New()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
