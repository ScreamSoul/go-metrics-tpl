package routers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gojuno/minimock/v3"
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
	mc := minimock.NewController(t)

	mockDB := repositories.NewMetricStorageMock(mc)
	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(MetricRouter(mockDB))
	defer ts.Close()

	mockDB.AddMock.Set(
		func(m metric.Metric) {
		},
	)

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
	mc := minimock.NewController(t)

	mockDB := repositories.NewMetricStorageMock(mc)
	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(MetricRouter(mockDB))
	defer ts.Close()

	var testTable = []struct {
		name           string
		url            string
		method         string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "fail type",
			url:            "/value/fake_gauge/MetricGauge1",
			method:         "GET",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:   "value gauge ok",
			url:    "/value/gauge/MetricGauge1",
			method: "GET",
			mockSetup: func() {
				mockDB.GetMock.When(metric.Gauge, metric.MetricName("MetricGauge1")).Then("1.11", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "1.11",
		},
		{
			name:   "value counter ok",
			url:    "/value/counter/MetricCounter1",
			method: "GET",
			mockSetup: func() {
				mockDB.GetMock.When(metric.Counter, metric.MetricName("MetricCounter1")).Then("1", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "1",
		},
		{
			name:   "value counter ok",
			url:    "/value/counter/MetricCounter3",
			method: "GET",
			mockSetup: func() {
				mockDB.GetMock.When(metric.Counter, metric.MetricName("MetricCounter3")).Then("", errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:   "value gauge not found",
			url:    "/value/gauge/MetricGauge3",
			method: "GET",
			mockSetup: func() {
				mockDB.GetMock.When(metric.Gauge, metric.MetricName("MetricGauge3")).Then("", errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			v.mockSetup()
			resp := testRequest(t, ts, v.method, v.url)
			require.Equal(t, v.expectedStatus, resp.StatusCode())
			if v.expectedBody != "" {
				assert.Equal(t, v.expectedBody, string(resp.Body()))
			}
		})

	}

}

func TestListRouter(t *testing.T) {
	mc := minimock.NewController(t)

	mockMetricList := []metric.Metric{
		{Name: "MetricGauge1", Value: "1.11", Type: metric.Gauge},
		{Name: "MetricGauge1", Value: "2.22", Type: metric.Gauge},
		{Name: "MetricCounter1", Value: "1", Type: metric.Counter},
		{Name: "MetricCounter1", Value: "2", Type: metric.Counter},
	}

	mockDB := repositories.NewMetricStorageMock(mc).ListMock.Return(
		mockMetricList,
	)

	defer mockDB.MinimockFinish()

	ts := httptest.NewServer(
		MetricRouter(
			mockDB,
		),
	)
	defer ts.Close()

	var testTable = []struct {
		name   string
		url    string
		method string
		status int
		result []metric.Metric
	}{
		{"list json ok", "/", "GET", http.StatusOK, mockMetricList},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp := testRequest(t, ts, v.method, v.url)
			require.Equal(t, v.status, resp.StatusCode())
			if v.status == http.StatusOK {
				var mList []metric.Metric
				err := json.Unmarshal(resp.Body(), &mList)
				require.NoError(t, err)
				assert.Equal(t, v.result, mList)
			}
		})

	}

}
