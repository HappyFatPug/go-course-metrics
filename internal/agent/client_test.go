package agent

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostMetrics(t *testing.T) {
	var gaugeValue = 0.11
	var counterValue float64 = 11

	testCases := []struct {
		name        string
		m           *models.Metrics
		wantURLPath string
		wantError   bool
	}{
		{
			name: "counter post",
			m: &models.Metrics{
				ID:    "test counter",
				MType: models.Counter,
				Value: &counterValue,
			},
			wantURLPath: "/update/counter/test counter/11",
			wantError:   false,
		},
		{
			name: "gauge post",
			m: &models.Metrics{
				ID:    "test gauge",
				MType: models.Gauge,
				Value: &gaugeValue,
			},
			wantURLPath: "/update/gauge/test gauge/0.11",
			wantError:   false,
		},
		{
			name: "wrong type",
			m: &models.Metrics{
				ID:    "test counter",
				MType: "models.Counter",
				Value: &counterValue,
			},
			wantURLPath: "",
			wantError:   true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Exactly(t, test.wantURLPath, r.URL.Path)

				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			parsedTestURL, err := url.Parse(server.URL)
			require.NoError(t, err)

			port, err := strconv.ParseUint(parsedTestURL.Port(), 10, 64)
			require.NoError(t, err)

			client := Client{
				c:      server.Client(),
				scheme: parsedTestURL.Scheme,
				host:   parsedTestURL.Hostname(),
				port:   uint(port),
				logger: slog.Default(),
			}

			err = client.postMetrics(test.m)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
