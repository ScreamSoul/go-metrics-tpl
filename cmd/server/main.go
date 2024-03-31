package main

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/middlewares"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
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

	var mStorage repositories.MetricStorage
	if cfg.FileStoragePath != "" {
		mRestoreStorage := memory.NewRestoreMetricStorage(
			cfg.FileStoragePath,
			cfg.StoreInterval,
			cfg.Restore,
			logger.Log,
		)
		mStorage = mRestoreStorage
		defer mRestoreStorage.Save()
	} else {
		mStorage = memory.NewMemStorage()
	}

	var metricServer = handlers.NewMetricServer(
		mStorage,
		logger.Log,
	)

	var router = routers.NewMetricRouter(
		metricServer,
		middlewares.NewLoggingMiddleware(logger.Log).Middleware,
		middlewares.GzipRequestMiddleware,
		middlewares.GzipResponseMiddleware,
	)

	logger.Log.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}

}
