package metrics

import (
	"encoding/json"
	"fmt"
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
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
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
