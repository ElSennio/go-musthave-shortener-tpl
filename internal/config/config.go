package config

import (
	"flag"
)

var ConfigFlags struct {
	ServerAddress string
	BaseURL       string
}

func init() {

	flag.StringVar(&ConfigFlags.ServerAddress, "a", "localhost:8080", "HTTP server listen address")
	flag.StringVar(&ConfigFlags.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
}
