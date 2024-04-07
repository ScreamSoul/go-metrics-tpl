package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type MetricServer struct {
	store repositories.MetricStorage
}

func NewMetricServer(metricRepo repositories.MetricStorage) *MetricServer {
	return &MetricServer{store: metricRepo}
}

func (ms *MetricServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	var metricObj metrics.Metrics

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&metricObj); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		metriObjPoint, err := metrics.NewMetric(
			r.PathValue("metric_type"),
			r.PathValue("metric_name"),
			r.PathValue("metric_value"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metricObj = *metriObjPoint
	}

	if err := metricObj.ValidateValue(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ms.store.Add(metricObj)
	logger.Debug("add metric object", zap.Any("metricObj", metricObj))
}

func (ms *MetricServer) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	metricObj, err := metrics.NewMetric(
		r.PathValue("metric_type"),
		r.PathValue("metric_name"),
		"",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.store.Get(metricObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(metricObj.GetValue())); err != nil {
		logger.Error("Error writing response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ms *MetricServer) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metricObj metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricObj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := ms.store.Get(&metricObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&metricObj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ms *MetricServer) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := ms.store.List()

	w.Header().Set("Content-Type", "text/html")

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
