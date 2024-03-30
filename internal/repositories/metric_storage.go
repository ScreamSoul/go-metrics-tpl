package repositories

import (
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
)

type MetricStorage interface {
	Add(m metrics.Metrics)
	Get(m *metrics.Metrics) error
	List() []metrics.Metrics
}
