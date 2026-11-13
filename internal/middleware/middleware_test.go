package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	hfunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("RESPONSE"))
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(WithLogging(logger)(hfunc))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	type LogRecord struct {
		Time     string `json:"time"`
		Level    string `json:"level"`
		Msg      string `json:"msg"`
		URI      string `json:"uri"`
		Method   string `json:"method"`
		Status   int    `json:"status"`
		Duration int    `json:"duration"`
		Size     int    `json:"size"`
	}

	var logRecord LogRecord
	err = json.Unmarshal(buf.Bytes(), &logRecord)

	assert.NoError(t, err)

	assert.Equal(t, "INFO", logRecord.Level)
	assert.Equal(t, 8, logRecord.Size)
	assert.True(t, logRecord.Duration > 0)
}

func TestJSONHandlerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	hfunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(http.HandlerFunc(JSONHandler(hfunc, logger)))
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
