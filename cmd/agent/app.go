package main

import (
	"context"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/pkg/backoff"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

func sender(
	ctx context.Context,
	metricRepo *memory.CollectionMetricStorage,
	backoffIntervals []time.Duration,
	metricClient *client.MetricsClient,
	reportInterval time.Duration,
) {
	logger := logging.GetLogger()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricsList, err := metricRepo.List(ctx)
			if err != nil {
				panic(err)
			}

			sendMetric := func() error {
				return metricClient.SendMetric(ctx, metricsList)
			}
			if backoffIntervals != nil {
				if err := backoff.RetryWithBackoff(backoffIntervals, client.IsTemporaryNetworkError, sendMetric); err != nil {
					logger.Error("retry send metric error", zap.Error(err))
				}
			} else {
				if err := sendMetric(); err != nil {
					logger.Error("send metric error", zap.Error(err))
				}
			}
		}

		time.Sleep(reportInterval)
	}
}

func updater(
	ctx context.Context,
	metricRepo *memory.CollectionMetricStorage,
	pollInterval time.Duration,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricRepo.Update()
			metricRepo.UpdateRuntime()
			metricRepo.UpdateGopsutil()
			time.Sleep(pollInterval)
		}
	}
}
