package server

type Config struct {
	Address string
}

func DefaultConfig() Config {
	return Config{
		Address: ":7000",
	}
}
