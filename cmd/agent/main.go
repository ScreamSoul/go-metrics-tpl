package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/pkg/logger"
	"go.uber.org/zap"
)

func sendMetric(uploadURL string, metric metrics.Metrics) {
	jsonData, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(jsonData)
	if err != nil {
		panic(err)
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf.Bytes()).
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
