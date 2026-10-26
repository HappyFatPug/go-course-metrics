package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
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
