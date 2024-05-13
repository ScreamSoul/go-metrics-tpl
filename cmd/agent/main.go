package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logging.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger := logging.GetLogger()

	logger.Info("start agent")
	logger.Info("use metric server", zap.String("server", cfg.GetServerURL()))

	metricRepo := memory.NewCollectionMetricStorage()

	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	metricClient := client.NewMetricsClient(cfg.CompressRequest, cfg.HashBodyKey, cfg.GetUpdateMetricURL())

	go updater(ctx, metricRepo, pollInterval)
	logger.Info("start senders", zap.Int("count_senders", cfg.RateLimit))
	for i := 0; i < cfg.RateLimit; i++ {
		go sender(ctx, metricRepo, cfg.BackoffIntervals, metricClient, reportInterval)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Agent closed:", ctx.Err())
}
