package handlers

import (
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
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "", http.StatusMethodNotAllowed)
	// 	return
	// }

	var metricObj = metric.Metric{}

	metricObj.Type = metric.MetricType(r.PathValue("metric_type"))
	metricObj.Name = r.PathValue("metric_name")
	metricObj.Value = r.PathValue("metric_value")

	if !metricObj.Type.IsValid() || !metricObj.IsValidValue() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	(*ms.store).Add(metricObj)
	fmt.Println(metricObj)
	w.WriteHeader(200)
}
