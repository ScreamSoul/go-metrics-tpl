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

func (c *Config) GetServerURL() string {
	return fmt.Sprintf("http://%s", c.ListenServerHost)

}

func (c *Config) GetUpdateMetricURL(mType, mName, mValue string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", c.GetServerURL(), mType, mName, mValue)
}

func NewConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.ListenServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "the frequency of sending metrics to the server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the runtime package")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
