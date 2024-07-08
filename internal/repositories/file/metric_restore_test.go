package file_test

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Successfully opens or creates the file for writing
func TestFileRestoreMetricWrapper_Save_Success(t *testing.T) {
	ctrl := minimock.NewController(t)

	mockMetricService := NewMetricStorageMock(ctrl)

	ctx := context.Background()

	fileTemp, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		errRemove := os.Remove(fileTemp.Name())
		assert.NoError(t, errRemove)
	}()

	wrapper := file.NewFileRestoreMetricWrapper(
		ctx, mockMetricService, fileTemp.Name(), 1, false,
	)

	metricsList := []metrics.Metrics{{ID: "test_metric", MType: metrics.Gauge, Value: new(float64)}}

	mockMetricService.ListMock.Return(metricsList, nil)

	wrapper.Save(ctx)

	fileContent, err := os.ReadFile(fileTemp.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	var savedMetrics []metrics.Metrics
	if err := json.Unmarshal(fileContent, &savedMetrics); err != nil {
		t.Fatalf("failed to unmarshal saved metrics: %v", err)
	}

	if !reflect.DeepEqual(savedMetrics, metricsList) {
		t.Errorf("expected %v, got %v", metricsList, savedMetrics)
	}
}

func TestFileRestoreMetricWrapper_Load_Success(t *testing.T) {
	ctrl := minimock.NewController(t)

	mockMetricService := NewMetricStorageMock(ctrl)

	ctx := context.Background()

	fileTemp, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		errRemove := os.Remove(fileTemp.Name())
		assert.NoError(t, errRemove)
	}()

	wrapper := file.NewFileRestoreMetricWrapper(
		ctx, mockMetricService, fileTemp.Name(), 0, false,
	)

	metricsData := []metrics.Metrics{{ID: "test_metric", MType: metrics.Gauge, Value: new(float64)}}
	err = json.NewEncoder(fileTemp).Encode(metricsData)
	require.NoError(t, err)
	err = fileTemp.Close()
	assert.NoError(t, err)

	mockMetricService.BulkAddMock.Expect(ctx, metricsData).Return(nil)

	wrapper.Load(ctx)
}

func TestGetMetricSuccessfully(t *testing.T) {
	ctrl := minimock.NewController(t)

	mockMetricService := NewMetricStorageMock(ctrl)

	ctx := context.Background()

	wrapper := file.NewFileRestoreMetricWrapper(
		ctx, mockMetricService, "", 0, false,
	)

	metric := &metrics.Metrics{ID: "test_metric"}

	mockMetricService.GetMock.Expect(ctx, metric).Return(nil)

	err := wrapper.Get(ctx, metric)

	assert.NoError(t, err)
}

func TestListMetricSuccessfully(t *testing.T) {
	ctrl := minimock.NewController(t)

	mockMetricService := NewMetricStorageMock(ctrl)

	ctx := context.Background()

	wrapper := file.NewFileRestoreMetricWrapper(
		ctx, mockMetricService, "", 0, false,
	)

	metricsList := []metrics.Metrics{}

	mockMetricService.ListMock.Expect(ctx).Return(metricsList, nil)

	_, err := wrapper.List(ctx)

	assert.NoError(t, err)
}
