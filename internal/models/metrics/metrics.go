package metrics

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type MetricType string

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

type Metrics struct {
	ID    string     `json:"id" db:"name"`                         // имя метрики
	MType MetricType `json:"type" db:"m_type"`                     // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty" db:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty" db:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetric(metricType, metricName, metricValue string) (*Metrics, error) {
	mType := MetricType(metricType)
	if !mType.IsValid() {
		return nil, fmt.Errorf("invalid metric type: %s", metricType)
	}

	metrics := &Metrics{
		ID:    metricName,
		MType: mType,
	}

	if metricValue == "" {
		return metrics, nil
	}

	switch mType {
	case Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric value as float64: %w", err)
		}
		metrics.Value = &value
	case Counter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse metric delta as int64: %w", err)
		}
		metrics.Delta = &delta
	}

	return metrics, nil
}

func (m *Metrics) GetValue() (val string) {
	switch m.MType {
	case Gauge:
		val = fmt.Sprint(*m.Value)
	case Counter:
		val = fmt.Sprint(*m.Delta)
	}
	return
}

func (m *Metrics) ValidateType() error {
	if !m.MType.IsValid() {
		return fmt.Errorf("metric type `%v` is not valid", m.MType)
	}
	return nil
}

func (m *Metrics) ValidateValue() error {
	switch m.MType {
	case Gauge:
		if m.Value == nil {
			return fmt.Errorf("metric type `%s` must be set Value filed", m.MType)
		}
	case Counter:
		if m.Delta == nil {
			return fmt.Errorf("metric type `%s` must be set Delta filed", m.MType)
		}
	}
	return nil
}

func (m *Metrics) UnmarshalJSON(data []byte) error {
	type Alias Metrics
	aux := &struct {
		MType string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.MType = MetricType(aux.MType)

	return m.ValidateType()
}
