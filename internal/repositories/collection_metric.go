package repositories

import (
	"context"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
)

//go:generate minimock -i github.com/screamsoul/go-metrics-tpl/internal/repositories.CollectionMetric -o ./mocks/collection_metric_mock.go -g
type CollectionMetric interface {
	Update()
	UpdateRuntime()
	UpdateGopsutil()
	List(ctx context.Context) ([]metrics.Metrics, error)
}
