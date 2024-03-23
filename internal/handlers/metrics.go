package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"go.uber.org/zap"
)

type MetricServer struct {
	store  repositories.MetricStorage
	logger *zap.Logger
}

func NewMetricServer(metricRepo repositories.MetricStorage, logger *zap.Logger) *MetricServer {
	return &MetricServer{store: metricRepo, logger: logger}
}

func (ms *MetricServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricObj, err := metric.NewMetric(
		r.PathValue("metric_type"),
		r.PathValue("metric_name"),
		r.PathValue("metric_value"),
	)

	if err != nil || !metricObj.IsValidValue() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ms.store.Add(metricObj)
	ms.logger.Debug("add metric object", zap.Any("metricObj", metricObj))
}

func (ms *MetricServer) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	mt := metric.MetricType(r.PathValue("metric_type"))
	mn := metric.MetricName(r.PathValue("metric_name"))

	if !mt.IsValid() {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	mv, err := ms.store.Get(mt, mn)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(mv)); err != nil {
		ms.logger.Error("Error writing response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ms *MetricServer) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := ms.store.List()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
