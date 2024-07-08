// correctly compresses a valid byte slice body
package client_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/client"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/stretchr/testify/assert"
)

func TestNewGzipCompressBodyMiddleware_CompressesValidBody(t *testing.T) {
	middleware := client.NewGzipCompressBodyMiddleware()
	req := resty.New().R()
	req.SetBody([]byte("test body"))

	err := middleware(nil, req)
	assert.NoError(t, err)

	compressedBody := req.Body.([]byte)
	buf := bytes.NewBuffer(compressedBody)
	gz, err := gzip.NewReader(buf)
	assert.NoError(t, err)

	decompressedBody, err := io.ReadAll(gz)
	assert.NoError(t, err)
	assert.Equal(t, "test body", string(decompressedBody))
}

func TestCorrectlySetsHashSHA256Header(t *testing.T) {
	hashKey := "testKey"
	body := []byte("testBody")
	expectedHash := sha256.New()
	expectedHash.Write(body)
	expectedHash.Write([]byte(hashKey))
	expectedHashSum := expectedHash.Sum(nil)

	middleware := client.NewHashSumHeaderMiddleware(hashKey)
	req := resty.New().R().SetBody(body)

	err := middleware(nil, req)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedHashSum), req.Header.Get("HashSHA256"))
}

// Successfully sends a list of metrics to the specified upload URL
func TestSendMetric_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a MetricsClient instance
	client := client.NewMetricsClient(
		false, "", server.URL,
	)

	// Create a context
	ctx := context.Background()

	// Create a sample metrics list
	metricsList := []metrics.Metrics{
		{ID: "metric1", MType: "gauge", Value: new(float64)},
		{ID: "metric2", MType: "counter", Delta: new(int64)},
	}

	// Call the SendMetric method
	err := client.SendMetric(ctx, metricsList)

	// Assert no error occurred
	assert.NoError(t, err)
}
