package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	ListenServerHost string `env:"ADDRESS"`
	ReportInterval   int    `env:"REPORT_INTERVAL"`
	PollInterval     int    `env:"POLL_INTERVAL"`
}

var cfg Config

func init() {
	flag.StringVar(&cfg.ListenServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "the frequency of sending metrics to the server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the runtime package")

}

func parseConfig() {
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Failed to parse env: %v", err)
	}

	fmt.Println(cfg)
}
