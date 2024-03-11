package routers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricRouter(t *testing.T) {

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
		resp, _ := testRequest(t, ts, v.method, v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		resp.Body.Close()
	}

}
