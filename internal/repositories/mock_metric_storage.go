package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
)

type MockMetricStorage struct {
	sync.Mutex
	Gauge   map[metric.MetricName]float64
	Counter map[metric.MetricName]int64
}

func NewMockMetricStorage() *MockMetricStorage {
	return &MockMetricStorage{
		Counter: make(map[metric.MetricName]int64),
		Gauge:   make(map[metric.MetricName]float64),
	}
}

func (db *MockMetricStorage) Add(m metric.Metric) {
	db.Lock()
	defer db.Unlock()

	switch m.Type {
	case metric.Gauge:
		val, _ := strconv.ParseFloat(string(m.Value), 64)
		db.Gauge[metric.MetricName(m.Name)] = val
	case metric.Counter:
		val, _ := strconv.ParseInt(string(m.Value), 10, 64)
		db.Counter[metric.MetricName(m.Name)] += val
	}
}

func (db *MockMetricStorage) Get(mt metric.MetricType, mn metric.MetricName) (string, error) {
	switch mt {
	case metric.Gauge:
		if v, ok := db.Gauge[mn]; ok {
			return fmt.Sprintf("%v", v), nil
		}
	case metric.Counter:
		if v, ok := db.Counter[mn]; ok {
			return fmt.Sprintf("%v", v), nil
		}
	}

	return "", errors.New("not fuund")
}

func (db *MockMetricStorage) List() (metics []metric.Metric) {
	for n, v := range db.Gauge {
		metics = append(metics, metric.Metric{
			Type:  metric.Gauge,
			Name:  n,
			Value: metric.MetricValue(strconv.FormatFloat(v, 'f', -1, 64)),
		})
	}
	for n, v := range db.Counter {
		metics = append(metics, metric.Metric{
			Type:  metric.Counter,
			Name:  n,
			Value: metric.MetricValue(fmt.Sprintf("%d", v)),
		})
	}
	return
}
