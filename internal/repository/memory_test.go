package repository

import (
	"slices"
	"testing"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestInMemMetricsRepositorySave(t *testing.T) {
	var counterValue float64 = 444
	var gaugeValue float64 = 445.33
	testCases := []struct {
		name      string
		m         *models.Metrics
		wantError bool
	}{
		{
			name:      "save counter ok",
			wantError: false,
			m:         &models.Metrics{ID: "1", MType: models.Counter, Value: &counterValue},
		},
		{
			name:      "save gauge ok",
			wantError: false,
			m:         &models.Metrics{ID: "1", MType: models.Gauge, Value: &gaugeValue},
		},
		{
			name:      "wrong type",
			wantError: true,
			m:         &models.Metrics{ID: "1", MType: "unknown", Value: &gaugeValue},
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
			var testCounterValue float64 = 11
			repo.Save(&models.Metrics{ID: "testCounter", MType: models.Counter, Value: &testCounterValue})
			var testGaugeValue float64 = 44.33
			repo.Save(&models.Metrics{ID: "testGauge", MType: models.Gauge, Value: &testGaugeValue})

			_, found := repo.Find(test.mType, test.mName)
			assert.Equal(t, test.wantFind, found)
		})
	}
}

func TestInMemMetricsRepositoryGetAll(t *testing.T) {
	repo := NewInMemMetricsRepository()

	var testCounterValue float64 = 11
	repo.Save(&models.Metrics{ID: "testCounter", MType: models.Counter, Value: &testCounterValue})
	var testGaugeValue float64 = 44.33
	repo.Save(&models.Metrics{ID: "testGauge", MType: models.Gauge, Value: &testGaugeValue})

	metrics := repo.GetAll()

	assert.Len(t, metrics, 2)

	for _, m := range metrics {
		assert.True(t, slices.Contains([]string{"testCounter", "testGauge"}, m.ID))
		assert.True(t, slices.Contains([]string{"11", "44.33"}, m.ValueStr()))
	}
}
