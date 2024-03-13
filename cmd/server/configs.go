package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ListenAddress string `env:"ADDRESS"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.ListenAddress, "a", "localhost:8080", "address and port to run server")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
