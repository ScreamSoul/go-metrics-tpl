package metric

import "strconv"

type MetricType string

type Metric struct {
	Type  MetricType
	Name  string
	Value string
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
		_, err := strconv.ParseFloat(mt.Value, 64)
		return err == nil
	case Counter:
		_, err := strconv.ParseInt(mt.Value, 10, 64)
		return err == nil
	default:
		return false
	}
}
