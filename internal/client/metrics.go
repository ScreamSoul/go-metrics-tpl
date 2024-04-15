package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type MetricsClient struct {
	resty.Client
	logger *zap.Logger
}

func GzipCompressBodyMiddleware(c *resty.Client, r *resty.Request) error {
	// Checking if there is already a Content-Encoding header
	if r.Header.Get("Content-Encoding") != "" {
		return nil
	}

	bodyBytes, ok := r.Body.([]byte)
	if !ok {
		return fmt.Errorf("body is not of type []byte")
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(bodyBytes); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	r.Body = buf.Bytes()

	r.Header.Set("Content-Encoding", "gzip")

	return nil
}

func NewMetricsClient(compressRequest bool) *MetricsClient {

	client := &MetricsClient{
		*resty.New(),
		logging.GetLogger(),
	}

	if compressRequest {
		client.OnBeforeRequest(GzipCompressBodyMiddleware)
	}

	return client
}

func (client *MetricsClient) SendMetric(ctx context.Context, uploadURL string, metricsList []metrics.Metrics) error {
	jsonData, err := json.Marshal(metricsList)
	if err != nil {
		panic(err)
	}

	resp, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(uploadURL)

	if err != nil {
		client.logger.Error("send error", zap.Error(err))
		return err
	}

	client.logger.Info(
		"send metric", zap.Any("metric", resp.Request.Body), zap.String("url", uploadURL),
	)
	return nil
}
