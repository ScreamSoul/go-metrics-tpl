package server_test

import (
	"context"
	"testing"

	"github.com/screamsoul/go-metrics-tpl/internal/server"
	"go.uber.org/zap"
)

func TestStart_InitializesInMemoryStorage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &server.Config{
		Postgres: server.Postgres{
			DatabaseDSN: "",
		},
		FileStoragePath: "",
		StoreInterval:   0,
		Restore:         false,
		HashBodyKey:     "",
		ListenAddress:   ":8080",
		Debug:           false,
	}
	logger := zap.NewNop()

	go server.Start(ctx, cfg, logger)
	cancel()
}
