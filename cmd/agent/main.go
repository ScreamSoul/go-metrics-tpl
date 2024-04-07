package main

import (
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
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

	logger.Info("start agent")
	logger.Info("use metric server", zap.String("server", cfg.GetServerURL()))

	metricRepo := memory.NewCollectionMetricStorage()

	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	metricClient := client.NewMetricsClient(cfg.CompressRequest)

	go func() {
		for {
			metricRepo.Update()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			for _, m := range metricRepo.List() {
				go metricClient.SendMetric(
					cfg.GetUpdateMetricURL(),
					m,
				)
			}
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
