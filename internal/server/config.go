package server

type Config struct {
	Address    string
	BaseDomain string
	Scheme     string
}

func DefaultConfig() Config {
	return Config{
		Address:    ":7000",
		BaseDomain: "",
		Scheme:     "https",
	}
}
