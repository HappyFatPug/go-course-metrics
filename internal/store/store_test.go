package store

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMetricsStore(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dump*.json")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	filePath := tmpFile.Name()
	tmpFile.Close()

	var intVal int64 = 100
	var floatVal float64 = 100.0

	testCases := []struct {
		name        string
		cfg         *config.StoreConfig
		testMetrics []models.Metrics
	}{
		{
			name: "basic test",
			cfg: &config.StoreConfig{
				StoreInterval:   1,
				FileStoragePath: filePath,
			},
			testMetrics: []models.Metrics{
				{
					ID:    "1",
					MType: models.Counter,
					Delta: &intVal,
				},
				{
					ID:    "2",
					MType: models.Gauge,
					Value: &floatVal,
				},
			},
		},

		{
			name: "basic test 1",
			cfg: &config.StoreConfig{
				StoreInterval:   0,
				FileStoragePath: filePath,
			},
			testMetrics: []models.Metrics{
				{
					ID:    "3",
					MType: models.Counter,
					Delta: &intVal,
				},
				{
					ID:    "4",
					MType: models.Gauge,
					Value: &floatVal,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			store := NewMetricsStore(test.cfg, slog.Default())

			assert.Equal(t, test.cfg.StoreInterval == 0, store.IsSync())

			assert.NoError(t, store.Dump(test.testMetrics))

			assert.FileExists(t, filePath)

			file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
			assert.NoError(t, err)

			var dumpedMetrics []models.Metrics
			testDec := json.NewDecoder(file)

			for {
				var m models.Metrics
				err := testDec.Decode(&m)
				if err != nil {
					assert.ErrorContains(t, err, "EOF")
					break
				}
				dumpedMetrics = append(dumpedMetrics, m)
			}

			file.Close()

			storedMetrics, err := store.ListStored()
			assert.NoError(t, err)
			assert.ObjectsAreEqual(storedMetrics, dumpedMetrics)
		})
	}
}
