package main

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/middlewares"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
	"github.com/screamsoul/go-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	var router = routers.NewMetricRouter(
		memory.NewMemStorage(),
		logger.Log,

		middlewares.NewLoggingMiddleware(logger.Log).Middleware,
	)

	logger.Log.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}
}
