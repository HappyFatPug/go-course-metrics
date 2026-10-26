package repository

import (
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type MetricsRepository interface {
	Find(mType, name string) (*models.Metrics, bool)
	Save(m *models.Metrics) error
}
