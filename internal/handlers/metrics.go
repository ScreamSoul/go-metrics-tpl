package handlers

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
)

type metricServer struct {
	store *repositories.MetricStorage
}

func NewMetricServer(metric_repo repositories.MetricStorage) *metricServer {
	return &metricServer{store: &metric_repo}
}

func (ms *metricServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var metric_obj = metric.Metric{}

	metric_obj.Type = metric.MetricType(r.PathValue("metric_type"))
	metric_obj.Name = r.PathValue("metric_name")
	metric_obj.Value = r.PathValue("metric_value")

	if !metric_obj.Type.IsValid() || !metric_obj.IsValidValue() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	(*ms.store).Add(metric_obj)

	w.WriteHeader(200)
}
