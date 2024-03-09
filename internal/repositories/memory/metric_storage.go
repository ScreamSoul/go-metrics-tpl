package memory

import (
	"strconv"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
	}
}

func (db *MemStorage) Add(m metric.Metric) {
	switch m.Type {
	case metric.Gauge:
		val, _ := strconv.ParseFloat(m.Value, 64)
		db.gauge[m.Name] = val
	case metric.Counter:
		val, _ := strconv.ParseInt(m.Value, 10, 64)
		db.counter[m.Name] += val
	}
}
