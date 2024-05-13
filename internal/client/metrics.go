package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type MetricsClient struct {
	resty.Client
	logger    *zap.Logger
	uploadURL string
}

func NewGzipCompressBodyMiddleware() func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
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

}

func NewHashSumHeaderMiddleware(hashKey string) func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		bodyBytes, ok := r.Body.([]byte)
		if !ok {
			return fmt.Errorf("body is not of type []byte")
		}
		h := sha256.New()

		h.Write(bodyBytes)
		h.Write([]byte(hashKey))

		dst := h.Sum(nil)

		r.Header.Set("HashSHA256", string(dst))

		return nil
	}
}

func NewMetricsClient(compressRequest bool, hashKey string, uploadURL string) *MetricsClient {

	client := &MetricsClient{
		*resty.New(),
		logging.GetLogger(),
		uploadURL,
	}

	if compressRequest {
		client.OnBeforeRequest(NewGzipCompressBodyMiddleware())
	}

	if hashKey != "" {
		client.OnBeforeRequest(NewHashSumHeaderMiddleware(hashKey))
	}

	return client
}

func (client *MetricsClient) SendMetric(ctx context.Context, metricsList []metrics.Metrics) error {
	jsonData, err := json.Marshal(metricsList)
	if err != nil {
		panic(err)
	}

	resp, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(client.uploadURL)

	if err != nil {
		client.logger.Error("send error", zap.Error(err))
		return err
	}

	client.logger.Info(
		"send metric", zap.Any("metric", resp.Request.Body), zap.String("url", client.uploadURL),
	)
	return nil
}
