package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"go.uber.org/zap"
)

func NewMetricRouter(
	storage repositories.MetricStorage,
	logger *zap.Logger,
	middlewares ...func(http.Handler) http.Handler,
) chi.Router {
	var metricServer = handlers.NewMetricServer(
		storage,
		logger,
	)
	r := chi.NewRouter()

	r.Use(middlewares...)

	r.Get("/", metricServer.ListMetrics)
	r.Post("/value/", metricServer.GetMetricJSON)
	r.Get("/value/{metric_type}/{metric_name}", metricServer.GetMetricValue)
	r.Post("/update/", metricServer.UpdateMetric)
	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", metricServer.UpdateMetric)

	return r
}
