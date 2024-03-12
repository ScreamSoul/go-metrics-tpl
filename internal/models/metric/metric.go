package metric

import (
	"fmt"
	"strconv"
)

type MetricType string
type MetricName string
type MetricValue string

type Metric struct {
	Type  MetricType  `json:"type"`
	Name  MetricName  `json:"name"`
	Value MetricValue `json:"value"`
}

func NewMetric(mType string, mName string, mValue string) (Metric, error) {

	if !MetricType(mType).IsValid() {
		return Metric{}, fmt.Errorf("metric type `%s` not valid", mType)
	}

	return Metric{
		Type:  MetricType(mType),
		Name:  MetricName(mName),
		Value: MetricValue(mValue),
	}, nil
}

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

func (mt MetricType) IsValid() bool {
	if mt == Gauge || mt == Counter {
		return true
	}
	return false

}

func (mt Metric) IsValidValue() bool {
	switch mt.Type {
	case Gauge:
		_, err := strconv.ParseFloat(string(mt.Value), 64)
		return err == nil
	case Counter:
		_, err := strconv.ParseInt(string(mt.Value), 10, 64)
		return err == nil
	default:
		return false
	}
}
