package server

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/middlewares"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/file"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/postgres"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
	"go.uber.org/zap"
)

// Start starts the server
func Start(ctx context.Context, cfg *Config, logger *zap.Logger) {
	// Create MetricStorage.
	var mStorage repositories.MetricStorage

	if cfg.DatabaseDSN == "" {
		// if no connection to the database is specified, the in-memory storage will be used.

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

	// Create restore wrapper.
	mStorageRestore := file.NewFileRestoreMetricWrapper(
		ctx,
		mStorage,
		cfg.FileStoragePath,
		cfg.StoreInterval,
		cfg.Restore,
	)

	if mStorageRestore.IsActiveRestore {
		defer mStorageRestore.Save(ctx)
	}

	var metricServer = handlers.NewMetricServer(
		mStorageRestore,
	)

	var router = routers.NewMetricRouter(
		metricServer,
		middlewares.LoggingMiddleware,
		middlewares.NewDecryptMiddleware(cfg.CryptoKey.Key),
		middlewares.NewHashSumHeaderMiddleware(cfg.HashBodyKey),
		middlewares.GzipDecompressMiddleware,
		middlewares.GzipCompressMiddleware,
	)

	if cfg.Debug {
		router.Mount("/debug", http.DefaultServeMux)
		logger.Info("mount debug pprof")
	}

	logger.Info("starting server", zap.String("ListenAddress", cfg.ListenAddress))

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}
}
