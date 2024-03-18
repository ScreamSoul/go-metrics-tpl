package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
	"github.com/screamsoul/go-metrics-tpl/internal/routers"
)

func main() {
	cfg, err := NewConfig()

	if err != nil {
		fmt.Printf("fail parse config: %v\r\n", err)
		os.Exit(1)
	}

	var router = routers.MetricRouter(
		memory.NewMemStorage(),
	)

	fmt.Println("Starting server on ", cfg.ListenAddress)

	if err := http.ListenAndServe(cfg.ListenAddress, router); err != nil {
		panic(err)
	}
}
