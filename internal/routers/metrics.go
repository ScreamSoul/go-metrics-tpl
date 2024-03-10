package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
)

func MetricRouter(storage repositories.MetricStorage) chi.Router {
	var metricServer = handlers.NewMetricServer(
		storage,
	)
	r := chi.NewRouter()

	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricServer.UpdateMetric)
	return r
}
