package repositories

import (
	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
)

type MetricStorage interface {
	Add(m metric.Metric)
	Get(mt metric.MetricType, mn metric.MetricName) (string, error)
	List() []metric.Metric
}
