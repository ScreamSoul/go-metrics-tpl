package main

import (
	"flag"
)

type Config struct {
	ListenHost string
}

var cfg Config

func init() {
	flag.StringVar(&cfg.ListenHost, "a", "localhost:8080", "address and port to run server")
}

func parseConfig() {
	flag.Parse()
}
