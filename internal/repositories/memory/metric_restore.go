package memory

import (
	"encoding/json"
	"os"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type RestoreMetricStorage struct {
	ms              MemStorage
	restoreFile     string
	restoreInterval int
	restoreInit     bool
}

func NewRestoreMetricStorage(
	restoreFile string,
	restoreInterval int,
	restoreInit bool,
) *RestoreMetricStorage {

	ms := &RestoreMetricStorage{
		*NewMemStorage(),
		restoreFile,
		restoreInterval,
		restoreInit,
	}

	if ms.restoreInit {
		ms.Load()
	}

	if ms.restoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(ms.restoreInterval) * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				ms.Save()
			}
		}()
	}

	return ms
}

func (db *RestoreMetricStorage) Save() {
	db.ms.Lock()
	defer db.ms.Unlock()
	logger := logging.GetLogger()

	logger.Info("Save metric to file")

	file, err := os.OpenFile(db.restoreFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("Error open or create file for write", zap.Error(err))
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(db.ms.List()); err != nil {
		logger.Error("Error saving metrics to file", zap.Error(err))
	}
}

func (db *RestoreMetricStorage) Load() {
	logger := logging.GetLogger()

	logger.Info("Load metric from file")

	file, err := os.OpenFile(db.restoreFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("Error open or create file for read", zap.Error(err))
		return
	}
	defer file.Close()

	fileInfo, err := os.Stat(db.restoreFile)
	if err != nil || fileInfo.Size() == 0 {
		logger.Warn("The open file has zero size or has just been created")
		return
	}

	metrics := []metrics.Metrics{}
	if err := json.NewDecoder(file).Decode(&metrics); err != nil {
		logger.Error("Error loading metrics from file", zap.Error(err))
		return
	}

	for _, metric := range metrics {
		db.ms.Add(metric)
	}
}

func (db *RestoreMetricStorage) Get(metric *metrics.Metrics) error {
	return db.ms.Get(metric)
}

func (db *RestoreMetricStorage) List() (metics []metrics.Metrics) {
	return db.ms.List()
}

func (db *RestoreMetricStorage) Add(m metrics.Metrics) {
	db.ms.Add(m)
	if db.restoreInterval == 0 {
		db.Save()
	}
}
