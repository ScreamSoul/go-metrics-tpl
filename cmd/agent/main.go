package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

func sendMetric(uploadURL string, body interface{}) {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(uploadURL)

	if err != nil {
		logger.Log.Error("send error", zap.Error(err))
	}

	if resp.IsError() {
		logger.Log.Error("error response", zap.Any("error", resp.Error()))
	}

	logger.Log.Info(
		"send metric", zap.Any("metric", resp.Request.Body), zap.String("url", uploadURL),
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
					cfg.GetUpdateMetricURL(),
					m,
				)
			}
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
