package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/astrahost/astrahost-tunnel/internal/client"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "http" {
		fmt.Printf("unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	port, err := strconv.ParseUint(os.Args[2], 10, 16)
	if err != nil || port == 0 {
		log.Fatalf("invalid port: %s", os.Args[2])
	}

	app := client.New(
		"127.0.0.1",
		uint16(port),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Println("Astra Tunnel")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  astra-tunnel http <port>")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  astra-tunnel http 5000")
	fmt.Println("  astra-tunnel http 3000")
	fmt.Println("  astra-tunnel http 8080")
}
