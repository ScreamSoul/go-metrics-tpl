package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type MemStorage struct {
	sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
	logger  *zap.Logger
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
		logger:  logging.GetLogger(),
	}
}

func (db *MemStorage) Add(ctx context.Context, m metrics.Metrics) error {
	db.Lock()
	defer db.Unlock()

	switch m.MType {
	case metrics.Gauge:
		db.gauge[m.ID] = *m.Value
	case metrics.Counter:
		db.counter[m.ID] += *m.Delta
	}
	return nil
}

func (db *MemStorage) Get(ctx context.Context, metric *metrics.Metrics) error {
	switch metric.MType {
	case metrics.Gauge:
		if v, ok := db.gauge[metric.ID]; ok {
			metric.Value = &v
			return nil
		}
	case metrics.Counter:
		if v, ok := db.counter[metric.ID]; ok {
			metric.Delta = &v
			return nil
		}
	}

	return errors.New("not found")
}

func (db *MemStorage) List(ctx context.Context) ([]metrics.Metrics, error) {
	metics := make([]metrics.Metrics, 0, len(db.counter)+len(db.gauge))
	for n, v := range db.gauge {
		metics = append(metics, metrics.Metrics{
			ID:    n,
			MType: metrics.Gauge,
			Value: &v,
		})
	}
	for n, v := range db.counter {
		metics = append(metics, metrics.Metrics{
			ID:    n,
			MType: metrics.Counter,
			Delta: &v,
		})
	}
	return metics, nil
}

func (db *MemStorage) Ping(ctx context.Context) bool {
	return true
}

func (db *MemStorage) BulkAdd(ctx context.Context, metricList []metrics.Metrics) error {
	for _, metric := range metricList {
		if err := db.Add(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}
