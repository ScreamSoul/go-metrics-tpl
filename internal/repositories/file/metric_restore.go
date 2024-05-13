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
	wrapper.logger.Info("save metric to file")

	file, err := os.OpenFile(wrapper.restoreFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		wrapper.logger.Error("error open or create file for write", zap.Error(err))
		return
	}
	defer file.Close()

	metricsList, err := wrapper.ms.List(ctx)
	if err != nil {
		wrapper.logger.Error("error read metric", zap.Error(err))
		return
	}

	if err := json.NewEncoder(file).Encode(metricsList); err != nil {
		wrapper.logger.Error("error saving metrics to file", zap.Error(err))
	}
}

func (wrapper *FileRestoreMetricWrapper) Load(ctx context.Context) {

	wrapper.logger.Info("load metric from file")

	file, err := os.OpenFile(wrapper.restoreFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		wrapper.logger.Error("error open or create file for read", zap.Error(err))
		return
	}
	defer file.Close()

	fileInfo, err := os.Stat(wrapper.restoreFile)
	if err != nil || fileInfo.Size() == 0 {
		wrapper.logger.Warn("the open file has zero size or has just been created")
		return
	}

	metrics := []metrics.Metrics{}
	if err := json.NewDecoder(file).Decode(&metrics); err != nil {
		wrapper.logger.Error("error loading metrics from file", zap.Error(err))
		return
	}
	err = wrapper.ms.BulkAdd(ctx, metrics)
	if err != nil {
		wrapper.logger.Error("error append metric to storage from file", zap.Error(err))
		return
	}
}

func (wrapper *FileRestoreMetricWrapper) Get(ctx context.Context, metric *metrics.Metrics) error {
	return wrapper.ms.Get(ctx, metric)
}

func (wrapper *FileRestoreMetricWrapper) List(ctx context.Context) ([]metrics.Metrics, error) {
	return wrapper.ms.List(ctx)
}

func (wrapper *FileRestoreMetricWrapper) Add(ctx context.Context, m metrics.Metrics) error {
	err := wrapper.ms.Add(ctx, m)

	if err != nil && wrapper.IsActiveRestore && wrapper.restoreInterval == 0 {
		wrapper.Save(ctx)
	}

	return err
}

func (wrapper *FileRestoreMetricWrapper) Ping(ctx context.Context) bool {
	return wrapper.ms.Ping(ctx)
}

func (wrapper *FileRestoreMetricWrapper) BulkAdd(ctx context.Context, metricList []metrics.Metrics) error {
	return wrapper.ms.BulkAdd(ctx, metricList)
}
