package main

import (
	"context"

	"github.com/screamsoul/go-metrics-tpl/internal/server"
	"github.com/screamsoul/go-metrics-tpl/internal/versions"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
)

func main() {
	versions.PrintBuildInfo()

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	cfg, err := server.NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logging.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger := logging.GetLogger()

	server.Start(ctx, cfg, logger)
}
