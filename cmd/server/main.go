package main

import (
	"fmt"
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
)

func main() {
	parseConfig()

	var router = routers.MetricRouter(
		memory.NewMemStorage(),
	)

	fmt.Println("Starting server on ", cfg.ListenAddress)

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}
}
