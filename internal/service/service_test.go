package service

import (
	"testing"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	testCases := []struct {
		name      string
		mType     string
		mName     string
		mValue    string
		wantError bool
	}{
		{
			name:      "simple gauge ok",
			mType:     models.Gauge,
			mName:     "name",
			mValue:    "55.4234523",
			wantError: false,
		},
		{
			name:      "parse gauge error",
			mType:     models.Gauge,
			mName:     "name",
			mValue:    "f55.4234523",
			wantError: true,
		},
		{
			name:      "simple counter ok",
			mType:     models.Counter,
			mName:     "name",
			mValue:    "55",
			wantError: false,
		},
		{
			name:      "parse counter error",
			mType:     models.Counter,
			mName:     "name",
			mValue:    "f55",
			wantError: true,
		},
		{
			name:      "wrong type",
			mType:     "counter1",
			mName:     "name",
			mValue:    "55",
			wantError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			logger.ConfigureLogger()

			if !test.wantError {
				assert.NoError(t, UpdateMetric(test.mType, test.mName, test.mValue))
			} else {
				assert.Error(t, UpdateMetric(test.mType, test.mName, test.mValue))
			}
		})
	}
}
