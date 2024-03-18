package memory

import (
	"runtime"
	"time"
)

type CollectionMetricStorage struct {
	MemStorage
}

func NewCollectionMetricStorage() *CollectionMetricStorage {
	return &CollectionMetricStorage{
		*NewMemStorage(),
	}
}

func (collection *CollectionMetricStorage) Update() {
	collection.Lock()
	defer collection.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	collection.gauge["Alloc"] = float64(mem.Alloc)
	collection.gauge["BuckHashSys"] = float64(mem.BuckHashSys)
	collection.gauge["Frees"] = float64(mem.Frees)
	collection.gauge["GCCPUFraction"] = mem.GCCPUFraction
	collection.gauge["GCSys"] = float64(mem.GCSys)
	collection.gauge["HeapAlloc"] = float64(mem.HeapAlloc)
	collection.gauge["HeapIdle"] = float64(mem.HeapIdle)
	collection.gauge["HeapInuse"] = float64(mem.HeapInuse)
	collection.gauge["HeapObjects"] = float64(mem.HeapObjects)
	collection.gauge["HeapReleased"] = float64(mem.HeapReleased)
	collection.gauge["HeapSys"] = float64(mem.HeapSys)
	collection.gauge["LastGC"] = float64(mem.LastGC)
	collection.gauge["Lookups"] = float64(mem.Lookups)
	collection.gauge["MCacheInuse"] = float64(mem.MCacheInuse)
	collection.gauge["MCacheSys"] = float64(mem.MCacheSys)
	collection.gauge["MSpanInuse"] = float64(mem.MSpanInuse)
	collection.gauge["MSpanSys"] = float64(mem.MSpanSys)
	collection.gauge["Mallocs"] = float64(mem.Mallocs)
	collection.gauge["NextGC"] = float64(mem.NextGC)
	collection.gauge["NumForcedGC"] = float64(mem.NumForcedGC)
	collection.gauge["NumGC"] = float64(mem.NumGC)
	collection.gauge["OtherSys"] = float64(mem.OtherSys)
	collection.gauge["PauseTotalNs"] = float64(mem.PauseTotalNs)
	collection.gauge["StackInuse"] = float64(mem.StackInuse)
	collection.gauge["StackSys"] = float64(mem.StackSys)
	collection.gauge["Sys"] = float64(mem.Sys)
	collection.gauge["TotalAlloc"] = float64(mem.TotalAlloc)

	collection.gauge["RandomValue"] = float64(time.Now().UnixNano()) / float64(time.Second)
	collection.counter["PollCount"]++
}
