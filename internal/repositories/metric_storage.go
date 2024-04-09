package repositories

import (
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
)

//go:generate minimock -i github.com/screamsoul/go-metrics-tpl/internal/repositories.MetricStorage -o ../mocks/metric_storage_mock.go -g
type MetricStorage interface {
	Add(m metrics.Metrics)
	Get(m *metrics.Metrics) error
	List() []metrics.Metrics
}
