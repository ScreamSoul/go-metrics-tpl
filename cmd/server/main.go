package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
)

func main() {
	flag.Parse()

	var router = routers.MetricRouter(
		memory.NewMemStorage(),
	)

	fmt.Println("Starting server on ", appFlags.listenHost)

	if err := http.ListenAndServe(appFlags.listenHost, router); err != nil {
		panic(err)
	}
}
