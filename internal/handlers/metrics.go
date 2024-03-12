package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
)

type MetricServer struct {
	store *repositories.MetricStorage
}

func NewMetricServer(metricRepo repositories.MetricStorage) *MetricServer {
	return &MetricServer{store: &metricRepo}
}

func (ms *MetricServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var metricObj, err = metric.NewMetric(
		r.PathValue("metric_type"),
		r.PathValue("metric_name"),
		r.PathValue("metric_value"),
	)

	if err != nil || !metricObj.IsValidValue() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	(*ms.store).Add(metricObj)
	fmt.Println(metricObj)
}

func (ms *MetricServer) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	var mt = metric.MetricType(r.PathValue("metric_type"))
	var mn = metric.MetricName(r.PathValue("metric_name"))

	if !mt.IsValid() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	mv, err := (*ms.store).Get(mt, mn)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	w.Write([]byte(mv))
}

func (ms *MetricServer) ListMetrics(w http.ResponseWriter, r *http.Request) {

	metrics := (*ms.store).List()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
