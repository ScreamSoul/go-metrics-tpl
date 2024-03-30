package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	ListenServerHost string `env:"ADDRESS"`
	ReportInterval   int    `env:"REPORT_INTERVAL"`
	PollInterval     int    `env:"POLL_INTERVAL"`
	LogLevel         string `env:"LOG_LEVEL"`
}

func (c *Config) GetServerURL() string {
	return strings.TrimRight(fmt.Sprintf("http://%s", c.ListenServerHost), "/")

}

func (c *Config) GetUpdateMetricURL() string {
	return fmt.Sprintf("%s/update/", c.GetServerURL())
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.ListenServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "the frequency of sending metrics to the server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the runtime package")
	flag.StringVar(&cfg.LogLevel, "ll", "INFO", "log level")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
