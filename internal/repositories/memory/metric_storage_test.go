package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/stretchr/testify/suite"
)

func newInt64(value int64) *int64 {
	v := new(int64)
	*v = value
	return v
}

func newFloat64(value float64) *float64 {
	v := new(float64)
	*v = value
	return v
}

type MemStorageSuite struct {
	suite.Suite
	storage *MemStorage
}

func TestMemStorageSuite(t *testing.T) {
	suite.Run(t, new(MemStorageSuite))
}

func (s *MemStorageSuite) SetupTest() {
	s.storage = NewMemStorage()
}

func (s *MemStorageSuite) TearDownTest() {
	s.storage.gauge = make(map[string]float64)
	s.storage.counter = make(map[string]int64)
}

func (s *MemStorageSuite) TestAdd() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		metric metrics.Metrics
		expect string
	}{
		{
			metric: metrics.Metrics{ID: "gauge1", MType: metrics.Gauge, Value: newFloat64(1.1)},
			expect: "1.1",
		},
		{
			metric: metrics.Metrics{ID: "counter1", MType: metrics.Counter, Delta: newInt64(1)},
			expect: "1",
		},
	}

	for _, tc := range testCases {
		err := s.storage.Add(ctx, tc.metric)
		s.Require().NoError(err)
		switch tc.metric.MType {
		case metrics.Gauge:
			s.Equal(tc.expect, fmt.Sprint(s.storage.gauge[tc.metric.ID]))
		case metrics.Counter:
			s.Equal(tc.expect, fmt.Sprint(s.storage.counter[tc.metric.ID]))
		}
	}
}

func (s *MemStorageSuite) TestGet() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		initDB func()
		metric metrics.Metrics
		expect string
	}{
		{
			initDB: func() {
				s.storage.gauge["gauge1"] = 1.1
			},
			metric: metrics.Metrics{ID: "gauge1", MType: metrics.Gauge},
			expect: "1.1",
		},
		{
			initDB: func() {
				s.storage.counter["counter1"] = 1
			},
			metric: metrics.Metrics{ID: "counter1", MType: metrics.Counter},
			expect: "1",
		},
	}

	for _, tc := range testCases {
		tc.initDB()

		err := s.storage.Get(ctx, &tc.metric)

		s.Require().NoError(err, s.storage)

		switch tc.metric.MType {
		case metrics.Gauge:
			s.Require().NotNil(tc.metric.Value)
			s.Equal(tc.expect, fmt.Sprint(*tc.metric.Value))
		case metrics.Counter:
			s.Require().NotNil(tc.metric.Delta)
			s.Equal(tc.expect, fmt.Sprint(*tc.metric.Delta))
		}
	}
}

func (s *MemStorageSuite) TestList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		initDB func()
		metric metrics.Metrics
		expect []metrics.Metrics
	}{
		{
			initDB: func() {
				s.storage.gauge["gauge1"] = 1.1
				s.storage.counter["counter1"] = 1
			},
			expect: []metrics.Metrics{
				{ID: "gauge1", MType: metrics.Gauge, Value: newFloat64(1.1)},
				{ID: "counter1", MType: metrics.Counter, Delta: newInt64(1)},
			},
		},
	}

	for _, tc := range testCases {
		tc.initDB()

		resMetrics, err := s.storage.List(ctx)
		s.NoError(err)
		s.Equal(tc.expect, resMetrics)

		s.TearDownTest()
	}
}
