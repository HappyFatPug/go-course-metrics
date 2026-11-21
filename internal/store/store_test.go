package store

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMetricsStore(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dump**.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(tmpFile.Name())
		fmt.Println(err)
	}()

	filePath := tmpFile.Name()
	tmpFile.Close()

	var intVal int64 = 100
	floatVal := 100.0

	testCases := []struct {
		name        string
		fpath       string
		testMetrics []models.Metrics
	}{
		{
			name:  "basic test",
			fpath: filePath,
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
			name:  "basic test 1",
			fpath: filePath,
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
			store := NewMetricsStore(test.fpath, slog.Default())

			assert.NoError(t, store.Dump(test.testMetrics))

			assert.FileExists(t, test.fpath)

			file, err := os.OpenFile(test.fpath, os.O_RDONLY, 0666)
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
			assert.ElementsMatch(t, test.testMetrics, storedMetrics)
			assert.ElementsMatch(t, storedMetrics, dumpedMetrics)
		})
	}
}
