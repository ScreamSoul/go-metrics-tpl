package repositories

import (
	"context"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
)

//go:generate minimock -i github.com/screamsoul/go-metrics-tpl/internal/repositories.MetricStorage -o ../mocks/metric_storage_mock.go -g
type MetricStorage interface {
	Add(ctx context.Context, m metrics.Metrics)
	Get(ctx context.Context, m *metrics.Metrics) error
	List(ctx context.Context) []metrics.Metrics
	Ping(ctx context.Context) bool
}
