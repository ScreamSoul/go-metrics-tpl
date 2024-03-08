package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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

type Storage interface {
	Add(metric Metric)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var storage = MemStorage{
	counter: make(map[string]int64),
	gauge:   make(map[string]float64),
}

func (db *MemStorage) Add(metric Metric) {
	switch metric.Type {
	case Gauge:
		val, _ := strconv.ParseFloat(metric.Value, 64)
		db.gauge[metric.Name] = val
	case Counter:
		val, _ := strconv.ParseInt(metric.Value, 10, 64)
		db.counter[metric.Name] += val
	}
}

func metricPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var metric = Metric{}

	metric.Type = MetricType(r.PathValue("metric_type"))
	metric.Name = r.PathValue("metric_name")
	metric.Value = r.PathValue("metric_value")

	if !metric.Type.IsValid() || !metric.IsValidValue() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	storage.Add(metric)

	fmt.Println("counter", storage.counter)
	fmt.Println("gauge", storage.gauge)

	w.WriteHeader(200)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", metricPage)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
