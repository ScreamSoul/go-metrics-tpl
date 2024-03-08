package main

import (
	"net/http"
)

type MetricType string
type MetricName string
type GaugeValue float64
type CounterValue int64

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

type MemStorage struct {
	db []string
}

type Storage interface {
	Add(metric_name string) bool
}

func metricPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var mt MetricType = MetricType(r.PathValue("metric_type"))
	// var mn MetricName = MetricName(r.PathValue("metric_name"))

	if !mt.IsValid() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", metricPage)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
