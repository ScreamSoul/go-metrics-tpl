package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
	"github.com/screamsoul/go-metrics-tpl/internal/utils"
)

type Metrics struct {
	sync.Mutex
	Alloc         uint64  `metric_type:"gauge"`
	BuckHashSys   uint64  `metric_type:"gauge"`
	Frees         uint64  `metric_type:"gauge"`
	GCCPUFraction float64 `metric_type:"gauge"`
	GCSys         uint64  `metric_type:"gauge"`
	HeapAlloc     uint64  `metric_type:"gauge"`
	HeapIdle      uint64  `metric_type:"gauge"`
	HeapInuse     uint64  `metric_type:"gauge"`
	HeapObjects   uint64  `metric_type:"gauge"`
	HeapReleased  uint64  `metric_type:"gauge"`
	HeapSys       uint64  `metric_type:"gauge"`
	LastGC        uint64  `metric_type:"gauge"`
	Lookups       uint64  `metric_type:"gauge"`
	MCacheInuse   uint64  `metric_type:"gauge"`
	MCacheSys     uint64  `metric_type:"gauge"`
	MSpanInuse    uint64  `metric_type:"gauge"`
	MSpanSys      uint64  `metric_type:"gauge"`
	Mallocs       uint64  `metric_type:"gauge"`
	NextGC        uint64  `metric_type:"gauge"`
	NumForcedGC   uint32  `metric_type:"gauge"`
	NumGC         uint32  `metric_type:"gauge"`
	OtherSys      uint64  `metric_type:"gauge"`
	PauseTotalNs  uint64  `metric_type:"gauge"`
	StackInuse    uint64  `metric_type:"gauge"`
	StackSys      uint64  `metric_type:"gauge"`
	Sys           uint64  `metric_type:"gauge"`
	TotalAlloc    uint64  `metric_type:"gauge"`
	PollCount     int64   `metric_type:"counter"`
	RandomValue   float64 `metric_type:"gauge"`
}

func (m *Metrics) updateMetrics() {
	m.Lock()
	defer m.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	m.Alloc = mem.Alloc
	m.BuckHashSys = mem.BuckHashSys
	m.Frees = mem.Frees
	m.GCCPUFraction = mem.GCCPUFraction
	m.GCSys = mem.GCSys
	m.HeapAlloc = mem.HeapAlloc
	m.HeapIdle = mem.HeapIdle
	m.HeapInuse = mem.HeapInuse
	m.HeapObjects = mem.HeapObjects
	m.HeapReleased = mem.HeapReleased
	m.HeapSys = mem.HeapSys
	m.LastGC = mem.LastGC
	m.Lookups = mem.Lookups
	m.MCacheInuse = mem.MCacheInuse
	m.MCacheSys = mem.MCacheSys
	m.MSpanInuse = mem.MSpanInuse
	m.MSpanSys = mem.MSpanSys
	m.Mallocs = mem.Mallocs
	m.NextGC = mem.NextGC
	m.NumForcedGC = mem.NumForcedGC
	m.NumGC = mem.NumGC
	m.OtherSys = mem.OtherSys
	m.PauseTotalNs = mem.PauseTotalNs
	m.StackInuse = mem.StackInuse
	m.StackSys = mem.StackSys
	m.Sys = mem.Sys
	m.TotalAlloc = mem.TotalAlloc
	m.PollCount++

	m.RandomValue = float64(time.Now().UnixNano()) / float64(time.Second)
}

func (m *Metrics) sendMetric(uploadURL string) {

	resp, err := http.Post(uploadURL, "text/plain", bytes.NewBufferString(""))
	if err != nil {
		fmt.Println(err)
		// panic(err)
		return
	}
	fmt.Printf("Url: %s; Status: %s\r\n", uploadURL, resp.Status)

	defer resp.Body.Close()
}

type ServerURL struct {
	baseURL string
}

func (su ServerURL) GetUpdateMetricURL(metric metric.Metric) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", strings.TrimRight(su.baseURL, "/"), metric.Type, metric.Name, metric.Value)
}

func main() {
	metrics := &Metrics{}

	pollInterval := time.Duration(appFlags.pollInterval) * time.Second
	reportInterval := time.Duration(appFlags.reportInterval) * time.Second
	serverURL := ServerURL{baseURL: fmt.Sprintf("http://%s/", appFlags.listenServerHost)}

	go func() {
		for {
			metrics.updateMetrics()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			fields, values := utils.PublicFields(metrics)
			for f := range fields {
				var value = fmt.Sprintf("%v", <-values)
				metric := metric.Metric{
					Name:  metric.MetricName(f.Name),
					Value: metric.MetricValue(value),
					Type:  metric.MetricType(f.Tag.Get("metric_type")),
				}
				go metrics.sendMetric(serverURL.GetUpdateMetricURL(metric))
			}
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
