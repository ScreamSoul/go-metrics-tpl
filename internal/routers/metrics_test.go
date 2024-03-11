package routers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metric"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) *resty.Response {

	req := resty.New().R()
	req.Method = method
	req.URL = ts.URL + path

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	return resp
}

func TestUpdateRouter(t *testing.T) {

	ts := httptest.NewServer(
		MetricRouter(
			repositories.NewMockMetricStorage(),
		),
	)
	defer ts.Close()

	var testTable = []struct {
		url    string
		method string
		status int
	}{
		{"/update/counter/someMetric/527", "POST", http.StatusOK},
		{"/update/gauge/someMetric1/527.44", "POST", http.StatusOK},
		{"/update/gauge/someMetric1/527.44", "GET", http.StatusMethodNotAllowed},
		{"/update/asd/someMetric1/527.44", "POST", http.StatusBadRequest},
		{"/update/gauge/someMetric1/asd", "POST", http.StatusBadRequest},
		{"/update/gauge", "POST", http.StatusNotFound},
	}
	for _, v := range testTable {
		resp := testRequest(t, ts, v.method, v.url)
		assert.Equal(t, v.status, resp.StatusCode())
	}

}

func TestValueRouter(t *testing.T) {
	mockDB := repositories.MockMetricStorage{
		Counter: map[metric.MetricName]int64{
			"MetricCounter1": 123,
			"MetricCounter2": 321,
		},
		Gauge: map[metric.MetricName]float64{
			"MetricGauge1": 1.12,
			"MetricGauge2": 1.32,
		},
	}

	ts := httptest.NewServer(
		MetricRouter(
			&mockDB,
		),
	)
	defer ts.Close()

	var testTable = []struct {
		name   string
		url    string
		method string
		status int
		text   string
	}{
		{"fail type", "/value/fake_gauge/MetricGauge1", "GET", http.StatusBadRequest, ""},
		{"value gauge ok", "/value/gauge/MetricGauge1", "GET", http.StatusOK, fmt.Sprintf("%v", mockDB.Gauge["MetricGauge1"])},
		{"value counter ok", "/value/counter/MetricCounter1", "GET", http.StatusOK, fmt.Sprintf("%v", mockDB.Counter["MetricCounter1"])},
		{"value counter not found", "/value/counter/MetricCounter3", "GET", http.StatusNotFound, ""},
		{"value gauge not found", "/value/gauge/MetricGauge3", "GET", http.StatusNotFound, ""},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, v.url)
			require.Equal(t, v.status, resp.StatusCode())
			if v.status == http.StatusOK {
				assert.Equal(t, v.text, string(resp.Body()))
			}
		})

	}

}

func TestListRouter(t *testing.T) {
	mockDB := repositories.MockMetricStorage{
		Counter: map[metric.MetricName]int64{
			"MetricCounter1": 123,
			"MetricCounter2": 321,
		},
		Gauge: map[metric.MetricName]float64{
			"MetricGauge1": 1.12,
			"MetricGauge2": 1.32,
		},
	}

	var respOkText = `[{"type":"gauge","name":"MetricGauge1","value":"1.12"},{"type":"gauge","name":"MetricGauge2","value":"1.32"},{"type":"counter","name":"MetricCounter1","value":"123"},{"type":"counter","name":"MetricCounter2","value":"321"}]`

	ts := httptest.NewServer(
		MetricRouter(
			&mockDB,
		),
	)
	defer ts.Close()

	var testTable = []struct {
		name   string
		url    string
		method string
		status int
		text   string
	}{
		{"list json ok", "/", "GET", http.StatusOK, respOkText},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, v.url)
			require.Equal(t, v.status, resp.StatusCode())
			if v.status == http.StatusOK {
				assert.JSONEq(t, string(resp.Body()), v.text)
			}
		})

	}

}
