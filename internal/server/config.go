package server

import "os"

type Config struct {
	Address    string
	BaseDomain string
	Scheme     string
}

func DefaultConfig() Config {
	address := os.Getenv("ASTRA_TUNNEL_ADDR")
	if address == "" {
		address = ":7000"
	}

	baseDomain := os.Getenv("ASTRA_BASE_DOMAIN")
	scheme := os.Getenv("ASTRA_SCHEME")
	if scheme == "" {
		scheme = "https"
	}

	return Config{
		Address:    address,
		BaseDomain: baseDomain,
		Scheme:     scheme,
	}
}
