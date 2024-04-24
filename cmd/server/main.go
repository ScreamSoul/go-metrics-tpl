package main

import (
	"context"
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
	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

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

	if cfg.DatabaseDSN == "" {
		memS := memory.NewMemStorage()
		mStorage = memS
	} else {
		postgresS := postgres.NewPostgresStorage(cfg.DatabaseDSN, cfg.BackoffIntervals)
		defer postgresS.Close()

		if err := postgresS.Bootstrap(ctx); err != nil {
			panic(err)
		}

		mStorage = postgresS
	}

	// Create restore wrapper
	mStorageRestore := file.NewFileRestoreMetricWrapper(
		ctx,
		mStorage,
		cfg.FileStoragePath,
		cfg.StoreInterval,
		cfg.Restore,
	)

	if mStorageRestore.IsActiveRestore {
		defer mStorageRestore.Save(context.Background())
	}

	var metricServer = handlers.NewMetricServer(
		mStorageRestore,
	)

	var router = routers.NewMetricRouter(
		metricServer,
		middlewares.LoggingMiddleware,
		middlewares.NewHashSumHeaderMiddleware(cfg.HashBodyKey),
		middlewares.GzipDecompressMiddleware,
		middlewares.GzipCompressMiddleware,
	)

	logger.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}

}
