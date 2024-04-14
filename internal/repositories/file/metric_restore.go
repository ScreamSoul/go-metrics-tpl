package file

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type FileRestoreMetricWrapper struct {
	ms              repositories.MetricStorage
	restoreFile     string
	restoreInterval int
	restoreInit     bool
	IsActiveRestore bool
	logger          *zap.Logger
}

func NewFileRestoreMetricWrapper(
	ctx context.Context,
	ms repositories.MetricStorage,
	restoreFile string,
	restoreInterval int,
	restoreInit bool,
) *FileRestoreMetricWrapper {

	restoreMetric := &FileRestoreMetricWrapper{
		ms:              ms,
		restoreFile:     restoreFile,
		restoreInterval: restoreInterval,
		restoreInit:     restoreInit,
		IsActiveRestore: restoreFile != "",
		logger:          logging.GetLogger(),
	}

	if restoreMetric.IsActiveRestore && restoreMetric.restoreInit {
		restoreMetric.Load(ctx)
	}

	if restoreMetric.IsActiveRestore && restoreMetric.restoreInterval > 0 {
		go func(context.Context) {
			ticker := time.NewTicker(time.Duration(restoreMetric.restoreInterval) * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				restoreMetric.Save(ctx)
			}
		}(ctx)
	}

	return restoreMetric
}

func (wrapper *FileRestoreMetricWrapper) Save(ctx context.Context) {
	wrapper.logger.Info("Save metric to file")

	file, err := os.OpenFile(wrapper.restoreFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		wrapper.logger.Error("Error open or create file for write", zap.Error(err))
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(wrapper.ms.List(ctx)); err != nil {
		wrapper.logger.Error("Error saving metrics to file", zap.Error(err))
	}
}

func (wrapper *FileRestoreMetricWrapper) Load(ctx context.Context) {

	wrapper.logger.Info("Load metric from file")

	file, err := os.OpenFile(wrapper.restoreFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		wrapper.logger.Error("Error open or create file for read", zap.Error(err))
		return
	}
	defer file.Close()

	fileInfo, err := os.Stat(wrapper.restoreFile)
	if err != nil || fileInfo.Size() == 0 {
		wrapper.logger.Warn("The open file has zero size or has just been created")
		return
	}

	metrics := []metrics.Metrics{}
	if err := json.NewDecoder(file).Decode(&metrics); err != nil {
		wrapper.logger.Error("Error loading metrics from file", zap.Error(err))
		return
	}

	for _, metric := range metrics {
		wrapper.ms.Add(ctx, metric)
	}
}

func (wrapper *FileRestoreMetricWrapper) Get(ctx context.Context, metric *metrics.Metrics) error {
	return wrapper.ms.Get(ctx, metric)
}

func (wrapper *FileRestoreMetricWrapper) List(ctx context.Context) (metics []metrics.Metrics) {
	return wrapper.ms.List(ctx)
}

func (wrapper *FileRestoreMetricWrapper) Add(ctx context.Context, m metrics.Metrics) {
	wrapper.ms.Add(ctx, m)
	if wrapper.IsActiveRestore && wrapper.restoreInterval == 0 {
		wrapper.Save(ctx)
	}
}

func (wrapper *FileRestoreMetricWrapper) Ping(ctx context.Context) bool {
	return wrapper.ms.Ping(ctx)
}
