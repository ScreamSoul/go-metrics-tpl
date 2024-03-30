package memory

import (
	"errors"
	"sync"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
)

type MemStorage struct {
	sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
	}
}

func (db *MemStorage) Add(m metrics.Metrics) {
	db.Lock()
	defer db.Unlock()

	switch m.MType {
	case metrics.Gauge:
		db.gauge[m.ID] = *m.Value
	case metrics.Counter:
		db.counter[m.ID] += *m.Delta
	}
}

func (db *MemStorage) Get(metric *metrics.Metrics) error {
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

	return errors.New("not fuund")
}

func (db *MemStorage) List() (metics []metrics.Metrics) {
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
	return
}
