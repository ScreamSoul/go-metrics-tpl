package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

func sendMetric(uploadURL string) {
	resp, err := http.Post(uploadURL, "text/plain", bytes.NewBufferString(""))
	if err != nil {
		logger.Log.Error("send error", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	logger.Log.Info(
		"send metric",
		zap.String("url", uploadURL),
		zap.String("resp status", resp.Status),
	)
}

func main() {
	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	logger.Log.Info("start agent")
	logger.Log.Info("use metric server", zap.String("server", cfg.GetServerURL()))

	metricRepo := memory.NewCollectionMetricStorage()

	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	go func() {
		for {
			metricRepo.Update()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			for _, m := range metricRepo.List() {
				go sendMetric(
					cfg.GetUpdateMetricURL(
						string(m.Type),
						string(m.Name),
						string(m.Value),
					),
				)
			}
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
