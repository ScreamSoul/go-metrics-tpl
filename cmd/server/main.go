package main

import (
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
)

func main() {
	var router = routers.MetricRouter(
		memory.NewMemStorage(),
	)
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		panic(err)
	}
}
