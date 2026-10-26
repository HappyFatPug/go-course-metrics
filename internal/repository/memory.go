package repository

import (
	"errors"
	"fmt"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type InMemMetricsRepository struct {
	gaugeBuf   map[string]float64
	counterBuf map[string]int64
}

func NewInMemMetricsRepository() *InMemMetricsRepository {
	return &InMemMetricsRepository{
		gaugeBuf:   make(map[string]float64),
		counterBuf: make(map[string]int64),
	}
}

func (repo *InMemMetricsRepository) Find(mType, name string) (*models.Metrics, bool) {
	switch mType {
	case models.Counter:
		value, ok := repo.counterBuf[name]
		if ok {
			return &models.Metrics{ID: name, MType: models.Counter, Delta: &value}, ok
		}
	case models.Gauge:
		value, ok := repo.gaugeBuf[name]
		if ok {
			return &models.Metrics{ID: name, MType: models.Gauge, Value: &value}, ok
		}
	}
	return nil, false
}

func (repo *InMemMetricsRepository) Save(m *models.Metrics) error {
	switch m.MType {
	case models.Counter:
		if m.Delta == nil {
			return errors.New("counter delta is undefined")
		}
		repo.counterBuf[m.ID] = *m.Delta
		logger.Logger.Debug(fmt.Sprintf("Counter saved: %s=%d", m.ID, *m.Delta))

	case models.Gauge:
		if m.Value == nil {
			return errors.New("gauge value is undefined")
		}
		repo.gaugeBuf[m.ID] = *m.Value
		logger.Logger.Debug(fmt.Sprintf("Gauge saved: %s=%.2f", m.ID, *m.Value))

	default:
		return errors.New("wrong type")
	}
	return nil
}
