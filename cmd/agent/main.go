package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories/memory"
)

func sendMetric(uploadURL string) {
	resp, err := http.Post(uploadURL, "text/plain", bytes.NewBufferString(""))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Url: %s; Status: %s\r\n", uploadURL, resp.Status)

	defer resp.Body.Close()
}

func main() {
	cfg, err := NewConfig()

	if err != nil {
		fmt.Println("fail parse config: ", err)
		os.Exit(1)
	}

	fmt.Print("start agent; ")
	fmt.Print("metric server: ", cfg.GetServerURL(), "; ")

	metricRepo := memory.NewCollectionMetricStorage()

	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	go func() {
		for {
			metricRepo.Update()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			for _, m := range metricRepo.List() {
				go sendMetric(
					cfg.GetUpdateMetricURL(
						string(m.Type),
						string(m.Name),
						string(m.Value),
					),
				)
			}
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
