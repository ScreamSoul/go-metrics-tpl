package memory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
)

type MemStorage struct {
	sync.Mutex
	gauge   map[metric.MetricName]float64
	counter map[metric.MetricName]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[metric.MetricName]int64),
		gauge:   make(map[metric.MetricName]float64),
	}
}

func (db *MemStorage) Add(m metric.Metric) {
	db.Lock()
	defer db.Unlock()

	switch m.Type {
	case metric.Gauge:
		val, _ := strconv.ParseFloat(string(m.Value), 64)
		db.gauge[metric.MetricName(m.Name)] = val
	case metric.Counter:
		val, _ := strconv.ParseInt(string(m.Value), 10, 64)
		db.counter[metric.MetricName(m.Name)] += val
	}
}

func (db *MemStorage) Get(mt metric.MetricType, mn metric.MetricName) (string, error) {
	switch mt {
	case metric.Gauge:
		if v, ok := db.gauge[mn]; ok {
			return fmt.Sprint(v), nil
		}
	case metric.Counter:
		if v, ok := db.counter[mn]; ok {
			return fmt.Sprint(v), nil
		}
	}

	return "", errors.New("not fuund")
}

func (db *MemStorage) List() (metics []metric.Metric) {
	for n, v := range db.gauge {
		metics = append(metics, metric.Metric{
			Type:  metric.Gauge,
			Name:  n,
			Value: metric.MetricValue(fmt.Sprint(v)),
		})
	}
	for n, v := range db.counter {
		metics = append(metics, metric.Metric{
			Type:  metric.Counter,
			Name:  n,
			Value: metric.MetricValue(fmt.Sprint(v)),
		})
	}
	return
}
