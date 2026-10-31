package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricHandler(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		url      string
		wantCode int
	}{
		{
			name:     "simple 200",
			method:   http.MethodPost,
			url:      "/update/counter/other/12",
			wantCode: http.StatusOK,
		},
		{
			name:     "404",
			method:   http.MethodPost,
			url:      "/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "405",
			method:   http.MethodPut,
			url:      "/update/counter/other/12",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "wrong type",
			method:   http.MethodPost,
			url:      "/update/1counter/other/12",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			app := app.ConfigureTestApp(false)

			router := mux.NewRouter()
			router.HandleFunc("/update/{type}/{name}/{value}", UpdateMetric(app)).Methods("POST")

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
			r := w.Result()

			assert.Equal(t, test.wantCode, r.StatusCode)
		})
	}
}

func TestGetMetric(t *testing.T) {
	type testMetric struct {
		ID    string
		mType string
		Value string
	}
	initialMetrics := []testMetric{
		{
			ID:    "testCounter",
			mType: models.Counter,
			Value: "123",
		},
		{
			ID:    "testGauge",
			mType: models.Gauge,
			Value: "0.321",
		},
	}
	testCases := []struct {
		name     string
		method   string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "counter 200",
			method:   http.MethodGet,
			url:      "/value/counter/testCounter",
			wantCode: http.StatusOK,
			wantBody: "123",
		},
		{
			name:     "gauge 200",
			method:   http.MethodGet,
			url:      "/value/gauge/testGauge",
			wantCode: http.StatusOK,
			wantBody: "0.321",
		},
		{
			name:     "counter 404",
			method:   http.MethodGet,
			url:      "/value/gauge/testGauge1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "counter 405",
			method:   http.MethodPost,
			url:      "/value/gauge/testGauge1",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "404",
			method:   http.MethodGet,
			url:      "/value/gauge1/test",
			wantCode: http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			app := app.ConfigureTestApp(false)

			for _, m := range initialMetrics {
				app.Service.UpdateMetric(m.mType, m.ID, m.Value)
			}

			router := mux.NewRouter()
			router.HandleFunc("/value/{type}/{name}", GetMetric(app)).Methods("GET")

			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
			r := w.Result()

			assert.Equal(t, test.wantCode, r.StatusCode)
			if test.wantCode == http.StatusOK {
				defer r.Body.Close()
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, test.wantBody, string(body))
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	type testMetric struct {
		ID    string
		mType string
		Value string
	}
	initialMetrics := []testMetric{
		{
			ID:    "testCounter",
			mType: models.Counter,
			Value: "123",
		},
		{
			ID:    "testGauge",
			mType: models.Gauge,
			Value: "0.321",
		},
	}

	app := app.ConfigureTestApp(false)

	for _, m := range initialMetrics {
		app.Service.UpdateMetric(m.mType, m.ID, m.Value)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", GetAll(app)).Methods("GET")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)
	r := w.Result()

	assert.Equal(t, http.StatusOK, r.StatusCode)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	bodyStr := string(body)
	bodyStr = strings.ReplaceAll(bodyStr, " ", "")
	// assert.Contains(t, bodyStr, "<td>testCounter</td>\n<td>123</td>")
	// assert.Contains(t, bodyStr, "<td>testGauge</td>\n<td>0.321</td>")
	assert.Contains(t, bodyStr, "<span>testCounter</span>:<span>123</span>")
	assert.Contains(t, bodyStr, "<span>testGauge</span>:<span>0.321</span>")
}
