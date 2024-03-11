package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

var cfg Config

type Config struct {
	ListenAddress string `env:"ADDRESS"`
}

func init() {
	flag.StringVar(&cfg.ListenAddress, "a", "localhost:8080", "address and port to run server")
}

func parseConfig() {
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Failed to parse: %v\r\n", err)
	}
}
