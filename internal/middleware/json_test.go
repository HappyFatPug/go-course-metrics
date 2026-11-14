package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONHandlerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	hfunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(http.HandlerFunc(JSONRequired(hfunc, logger)))
	defer server.Close()

	resp, err := server.Client().Post(server.URL, "text/plain", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	resp.Body.Close()

	resp, err = server.Client().Post(server.URL, "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
