package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	var tv int64 = 123
	testMetric := models.Metrics{
		ID:    "test metrics",
		MType: models.Counter,
		Delta: &tv,
	}
	testMetricBytes, err := json.Marshal(&testMetric)
	require.NoError(t, err)

	handlerF := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(testMetricBytes)
	}

	handler := GzipMiddleware(slog.Default())(http.HandlerFunc(handlerF))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("use Gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		gw := gzip.NewWriter(buf)
		_, err := gw.Write(testMetricBytes)
		require.NoError(t, err)
		err = gw.Close()
		require.NoError(t, err)

		r, err := http.NewRequest("POST", srv.URL, buf)
		require.NoError(t, err)
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := srv.Client().Do(r)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

		defer resp.Body.Close()

		reader, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		respBody, err := io.ReadAll(reader)
		require.JSONEq(t, string(testMetricBytes), string(respBody))
	})

	t.Run("no Gzip", func(t *testing.T) {
		r, err := http.NewRequest("POST", srv.URL, bytes.NewReader(testMetricBytes))
		require.NoError(t, err)

		resp, err := srv.Client().Do(r)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.JSONEq(t, string(testMetricBytes), string(respBody))
	})
}
