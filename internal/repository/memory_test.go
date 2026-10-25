package repository

import (
	"testing"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestInMemMetricsRepositorySave(t *testing.T) {
	var intVal int64 = 444
	var floatVal float64 = 445.33
	testCases := []struct {
		name      string
		m         *models.Metrics
		wantError bool
	}{
		{
			name:      "save counter ok",
			wantError: false,
			m:         &models.Metrics{ID: "1", MType: models.Counter, Delta: &intVal},
		},
		{
			name:      "save gauge ok",
			wantError: false,
			m:         &models.Metrics{ID: "1", MType: models.Gauge, Value: &floatVal},
		},
		{
			name:      "wrong type",
			wantError: true,
			m:         &models.Metrics{ID: "1", MType: "unknown", Value: &floatVal},
		},
		{
			name:      "wrong value",
			wantError: true,
			m:         &models.Metrics{ID: "1", MType: models.Counter, Value: &floatVal},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			repo := NewInMemMetricsRepository()

			if test.wantError {
				assert.Error(t, repo.Save(test.m))
			} else {
				assert.NoError(t, repo.Save(test.m))
			}
		})
	}
}

func TestInMemMetricsRepositoryFind(t *testing.T) {
	testCases := []struct {
		name     string
		mType    string
		mName    string
		wantFind bool
	}{
		{
			name:     "find ok",
			mType:    models.Counter,
			mName:    "testCounter",
			wantFind: true,
		},
		{
			name:     "find not ok",
			mType:    models.Gauge,
			mName:    "testCounter",
			wantFind: false,
		},
		{
			name:     "wrong type",
			mType:    "wrongType",
			mName:    "testCounter",
			wantFind: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			repo := NewInMemMetricsRepository()
			var testCounterValue int64 = 11
			repo.Save(&models.Metrics{ID: "testCounter", MType: models.Counter, Delta: &testCounterValue})
			var testGaugeValue float64 = 44.33
			repo.Save(&models.Metrics{ID: "testGauge", MType: models.Gauge, Value: &testGaugeValue})

			_, found := repo.Find(test.mType, test.mName)
			assert.Equal(t, test.wantFind, found)
		})
	}
}
