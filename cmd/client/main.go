package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/astrahost/astrahost-tunnel/internal/client"
)

const subdomainAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

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

	fs := flag.NewFlagSet("http", flag.ExitOnError)

	subdomain := fs.String("s", "", "subdomain")
	fs.StringVar(subdomain, "subdomain", "", "subdomain")

	if err := fs.Parse(os.Args[3:]); err != nil {
		log.Fatal(err)
	}

	if *subdomain == "" {
		*subdomain = randomSubdomain(8)
	}

	if err := validateSubdomain(*subdomain); err != nil {
		log.Fatal(err)
	}

	app := client.New(
		"127.0.0.1",
		uint16(port),
		*subdomain,
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func randomSubdomain(length int) string {
	result := make([]byte, length)
	max := big.NewInt(int64(len(subdomainAlphabet)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			log.Fatalf("failed to generate random subdomain: %v", err)
		}

		result[i] = subdomainAlphabet[n.Int64()]
	}

	return string(result)
}

func validateSubdomain(value string) error {
	if len(value) < 1 || len(value) > 63 {
		return fmt.Errorf("subdomain must be 1-63 characters")
	}

	for _, r := range value {
		if !strings.ContainsRune(subdomainAlphabet, r) && r != '-' {
			return fmt.Errorf("invalid subdomain character: %q", r)
		}
	}

	return nil
}

func printUsage() {
	fmt.Println("Astra Tunnel")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  astra-tunnel http <port>")
	fmt.Println("  astra-tunnel http <port> [-s <subdomain>]")
	fmt.Println("  astra-tunnel http <port> [--subdomain <subdomain>]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  astra-tunnel http 5000")
	fmt.Println("  astra-tunnel http 8080")
	fmt.Println("  astra-tunnel http 8080 -s custom")
	fmt.Println("  astra-tunnel http 8080 --subdomain custom")
}
