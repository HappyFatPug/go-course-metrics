package handler

import (
	"bytes"
	"encoding/json"
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
			router.HandleFunc("/update/{type}/{name}/{value}", UpdateMetric(app.Service, app.Logger)).Methods("POST")

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
				assert.NoError(t, app.Service.UpdateMetric(m.mType, m.ID, m.Value))
			}

			router := mux.NewRouter()
			router.HandleFunc("/value/{type}/{name}", GetMetric(app.Service)).Methods("GET")

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
		assert.NoError(t, app.Service.UpdateMetric(m.mType, m.ID, m.Value))
	}

	router := mux.NewRouter()
	router.HandleFunc("/", GetAll(app.Service, app.Logger)).Methods("GET")

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
	assert.Contains(t, bodyStr, "<span>testCounter</span>:<span>123</span>")
	assert.Contains(t, bodyStr, "<span>testGauge</span>:<span>0.321</span>")
}

func TestUpdateMetricJSONHandler(t *testing.T) {
	type postBodyMetrics struct {
		ID    string  `json:"id,omitempty"`
		MType string  `json:"type,omitempty"`
		Value float64 `json:"value,omitempty"`
	}
	testCases := []struct {
		name     string
		method   string
		url      string
		wantCode int
		m        postBodyMetrics
	}{
		{
			name:     "simple 200",
			method:   http.MethodPost,
			url:      "/update",
			wantCode: http.StatusOK,
			m: postBodyMetrics{
				ID:    "metric",
				MType: "counter",
				Value: 20,
			},
		},
		{
			name:     "404",
			method:   http.MethodPost,
			url:      "/update/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "405",
			method:   http.MethodPut,
			url:      "/update",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "validation error1",
			method:   http.MethodPost,
			url:      "/update",
			wantCode: http.StatusBadRequest,
			m: postBodyMetrics{
				ID:    "metric",
				MType: "counter1",
				Value: 20,
			},
		},
		{
			name:     "validation error2",
			method:   http.MethodPost,
			url:      "/update",
			wantCode: http.StatusBadRequest,
			m: postBodyMetrics{
				ID:    "metric",
				MType: "counter",
			},
		},
		{
			name:     "validation error3",
			method:   http.MethodPost,
			url:      "/update",
			wantCode: http.StatusBadRequest,
			m: postBodyMetrics{
				ID:    "metric",
				Value: 20,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			app := app.ConfigureTestApp(false)

			router := mux.NewRouter()
			router.HandleFunc("/update", UpdateMetricJSON(app.Service, app.Logger)).Methods("POST")

			body, err := json.Marshal(&test.m)
			assert.NoError(t, err)
			request := httptest.NewRequest(test.method, test.url, bytes.NewReader(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
			r := w.Result()

			assert.Equal(t, test.wantCode, r.StatusCode)
		})
	}
}

func TestGetMetricJSON(t *testing.T) {
	type responseBodyMetric struct {
		ID    string  `json:"id"`
		MType string  `json:"type"`
		Value float64 `json:"value"`
	}

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
		name      string
		method    string
		url       string
		wantCode  int
		mId       string
		mType     string
		wantValue float64
	}{
		{
			name:      "counter 200",
			method:    http.MethodGet,
			url:       "/value",
			wantCode:  http.StatusOK,
			mId:       "testCounter",
			mType:     models.Counter,
			wantValue: 123,
		},
		{
			name:      "gauge 200",
			method:    http.MethodGet,
			url:       "/value",
			wantCode:  http.StatusOK,
			mId:       "testGauge",
			mType:     models.Gauge,
			wantValue: 0.321,
		},
		{
			name:     "404",
			method:   http.MethodGet,
			url:      "/value/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "counter 404",
			method:   http.MethodGet,
			url:      "/value",
			wantCode: http.StatusNotFound,
			mId:      "testCounter1",
			mType:    models.Counter,
		},
		{
			name:     "405",
			method:   http.MethodPost,
			url:      "/value",
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			app := app.ConfigureTestApp(false)

			for _, m := range initialMetrics {
				assert.NoError(t, app.Service.UpdateMetric(m.mType, m.ID, m.Value))
			}

			router := mux.NewRouter()
			router.HandleFunc("/value", GetMetricJSON(app.Service, app.Logger)).Methods("GET")

			requestBodyM := struct {
				ID    string `json:"id"`
				MType string `json:"type"`
			}{
				ID:    test.mId,
				MType: test.mType,
			}
			body, err := json.Marshal(&requestBodyM)
			assert.NoError(t, err)
			request := httptest.NewRequest(test.method, test.url, bytes.NewReader(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)
			r := w.Result()

			assert.Equal(t, test.wantCode, r.StatusCode)
			if test.wantCode == http.StatusOK {
				defer r.Body.Close()
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)

				var responseM responseBodyMetric
				assert.NoError(t, json.Unmarshal(body, &responseM))
				assert.Equal(t, test.wantValue, responseM.Value)
			}
		})
	}
}
