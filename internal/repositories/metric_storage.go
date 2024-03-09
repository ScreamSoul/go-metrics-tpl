package repositories

import (
	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
)

type MetricStorage interface {
	Add(m metric.Metric)
}
