package main

import (
	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
)

func main() {

	cfg, err := client.NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logging.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger := logging.GetLogger()

	client.Start(cfg, logger)
}
