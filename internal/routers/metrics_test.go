package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/screamsoul/go-metrics-tpl/internal/handlers"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	body interface{},
	headers ...http.Header,
) *resty.Response {

	req := resty.New().R()

	req.Method = method
	req.URL = ts.URL + path

	if body != nil {
		req.Body = body
	}
	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				req.SetHeader(key, value)
			}
		}
	}

	resp, err := req.Send()
	require.NoError(t, err, "error making HTTP request")

	return resp
}

func TestUpdateRouter(t *testing.T) {
	mc := minimock.NewController(t)

	mockDB := repositories.NewMetricStorageMock(mc)
	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(
		NewMetricRouter(
			handlers.NewMetricServer(mockDB),
		),
	)
	defer ts.Close()

	mockDB.AddMock.Set(
		func(m metrics.Metrics) {
		},
	)

	var testTable = []struct {
		name   string
		body   map[string]interface{}
		method string
		status int
	}{
		{
			name:   "update ok",
			body:   map[string]interface{}{"type": "counter", "delta": 1, "ID": "someMetric1"},
			method: "POST",
			status: http.StatusOK,
		},
		{
			name:   "update mothod not allowed",
			body:   map[string]interface{}{"type": "counter", "delta": 1, "id": "someMetric1"},
			method: "GET",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "update bad type",
			body:   map[string]interface{}{"type": "asd", "delta": 1, "id": "someMetric1"},
			method: "POST",
			status: http.StatusBadRequest,
		},
		{
			name:   "update bad value",
			body:   map[string]interface{}{"type": "counter", "delta": "asd", "id": "someMetric1"},
			method: "POST",
			status: http.StatusBadRequest,
		},
		{
			name:   "update incorrect value",
			body:   map[string]interface{}{"type": "counter", "value": 1.2, "id": "someMetric1"},
			method: "POST",
			status: http.StatusBadRequest,
		},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, "/update/", v.body, http.Header{"Content-Type": {"application/json"}})
			assert.Equal(
				t,
				v.status,
				resp.StatusCode(),
				fmt.Sprintf("Resp body: %s", string(resp.Body())),
			)
		})
	}

}

func TestValueRouter(t *testing.T) {
	mc := minimock.NewController(t)

	mockDB := repositories.NewMetricStorageMock(mc)
	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(
		NewMetricRouter(
			handlers.NewMetricServer(mockDB),
		),
	)
	defer ts.Close()

	mockDB.GetMock.Set(func(m *metrics.Metrics) (err error) {
		var intValue int64 = 1
		if m.MType == metrics.Counter && m.ID == "someMetric1" {
			m.Delta = &intValue
		} else {
			return errors.New("not found")
		}

		return nil
	})

	var testTable = []struct {
		name           string
		body           map[string]interface{}
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "fail type",
			body:           map[string]interface{}{"type": "fake_type", "id": "someMetric1"},
			method:         "POST",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "value ok",
			body:           map[string]interface{}{"type": "counter", "id": "someMetric1"},
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"type": "counter", "id": "someMetric1", "delta": 1}`,
		},
		{
			name:   "value not found",
			body:   map[string]interface{}{"type": "counter", "id": "someMetric2"},
			method: "POST",

			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, "/value/", v.body, http.Header{"Content-Type": {"application/json"}})
			require.Equal(t, v.expectedStatus, resp.StatusCode(), fmt.Sprintf("Resp body: %s", string(resp.Body())))
			if v.expectedBody != "" {
				assert.JSONEq(t, v.expectedBody, string(resp.Body()))
			}
		})

	}

}

func TestListRouter(t *testing.T) {
	mc := minimock.NewController(t)

	var intValue int64 = 1
	var floatValue = 1.1

	mockMetricList := []metrics.Metrics{
		{ID: "MetricGauge1", Value: &floatValue, MType: metrics.Gauge},
		{ID: "MetricCounter1", Delta: &intValue, MType: metrics.Counter},
	}

	mockDB := repositories.NewMetricStorageMock(mc).ListMock.Return(
		mockMetricList,
	)

	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(
		NewMetricRouter(
			handlers.NewMetricServer(mockDB),
		),
	)
	defer ts.Close()

	var testTable = []struct {
		name   string
		method string
		status int
		result []metrics.Metrics
	}{
		{"list json ok", "GET", http.StatusOK, mockMetricList},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, "/", nil)
			require.Equal(t, v.status, resp.StatusCode())
			if v.status == http.StatusOK {
				var mList []metrics.Metrics
				err := json.Unmarshal(resp.Body(), &mList)
				require.NoError(t, err, string(resp.Body()))
				assert.Equal(t, v.result, mList)
			}
		})

	}

}
