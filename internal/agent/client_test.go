package agent

import (
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
	var counterValue int64 = 11
	var gaugeValue float64 = 0.11

	testCases := []struct {
		name        string
		m           *models.Metrics
		wantUrlPath string
		wantError   bool
	}{
		{
			name: "counter post",
			m: &models.Metrics{
				ID:    "test counter",
				MType: models.Counter,
				Delta: &counterValue,
			},
			wantUrlPath: "/update/counter/test counter/11",
			wantError:   false,
		},
		{
			name: "gauge post",
			m: &models.Metrics{
				ID:    "test gauge",
				MType: models.Gauge,
				Value: &gaugeValue,
			},
			wantUrlPath: "/update/gauge/test gauge/0.110000",
			wantError:   false,
		},
		{
			name: "wrong type",
			m: &models.Metrics{
				ID:    "test counter",
				MType: "models.Counter",
				Delta: &counterValue,
			},
			wantUrlPath: "/update/counter/test counter/11",
			wantError:   true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Exactly(t, test.wantUrlPath, r.URL.Path)

				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			parsedTestUrl, err := url.Parse(server.URL)
			require.NoError(t, err)

			port, err := strconv.ParseUint(parsedTestUrl.Port(), 10, 64)
			require.NoError(t, err)

			client := Client{
				c:      server.Client(),
				scheme: parsedTestUrl.Scheme,
				host:   parsedTestUrl.Hostname(),
				port:   uint(port),
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
