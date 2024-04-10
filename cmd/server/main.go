package main

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/middlewares"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/file"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/postgres"
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

	// Create MetricStorage

	var mStorage repositories.MetricStorage

	if cfg.DatabaseDNS == "" {
		memS := memory.NewMemStorage()
		mStorage = memS
	} else {
		postgresS := postgres.NewPostgresStorage(cfg.DatabaseDNS)
		defer postgresS.Close()
		mStorage = postgresS
	}

	// Create resore wrapper

	mStorageRestore := file.NewFileRestoreMetricWrapper(
		mStorage,
		cfg.FileStoragePath,
		cfg.StoreInterval,
		cfg.Restore,
	)

	if mStorageRestore.IsActiveRestore {
		defer mStorageRestore.Save()
	}

	var metricServer = handlers.NewMetricServer(
		mStorageRestore,
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
