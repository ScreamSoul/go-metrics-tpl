package main

import (
	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"

	"net/http"
)

func main() {
	mux := http.NewServeMux()

	var metricServer = handlers.NewMetricServer(
		memory.NewMemStorage(),
	)
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", metricServer.UpdateMetric)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
