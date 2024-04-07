package main

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/middlewares"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logging.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger := logging.GetLogger()

	var mStorage repositories.MetricStorage
	if cfg.FileStoragePath != "" {
		mRestoreStorage := memory.NewRestoreMetricStorage(
			cfg.FileStoragePath,
			cfg.StoreInterval,
			cfg.Restore,
		)
		mStorage = mRestoreStorage
		defer mRestoreStorage.Save()
	} else {
		mStorage = memory.NewMemStorage()
	}

	var metricServer = handlers.NewMetricServer(
		mStorage,
	)

	var router = routers.NewMetricRouter(
		metricServer,
		middlewares.LoggingMiddleware,
		middlewares.GzipDecompressMiddleware,
		middlewares.GzipCompressMiddleware,
	)

	logger.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}

}
