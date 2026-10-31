package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type InMemMetricsRepository struct {
	gaugeBuf   map[string]float64
	gaugeMu    sync.RWMutex
	counterBuf map[string]int64
	counterMu  sync.RWMutex
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
		repo.counterMu.RLock()
		defer repo.counterMu.RUnlock()
		value, ok := repo.counterBuf[name]
		valueFloat := float64(value)
		if ok {
			return &models.Metrics{ID: name, MType: models.Counter, Value: &valueFloat}, ok
		}
	case models.Gauge:
		repo.gaugeMu.RLock()
		defer repo.gaugeMu.RUnlock()
		value, ok := repo.gaugeBuf[name]
		if ok {
			return &models.Metrics{ID: name, MType: models.Gauge, Value: &value}, ok
		}
	}
	return nil, false
}

func (repo *InMemMetricsRepository) Save(m *models.Metrics) error {
	if m.Value == nil {
		return errors.New("value is undefined")
	}
	switch m.MType {
	case models.Counter:
		repo.counterMu.Lock()
		defer repo.counterMu.Unlock()
		repo.counterBuf[m.ID] = int64(*m.Value)
		logger.Logger.Debug(fmt.Sprintf("Counter saved: %s=%s", m.ID, m.ValueStr()))

	case models.Gauge:
		repo.gaugeMu.Lock()
		defer repo.gaugeMu.Unlock()
		repo.gaugeBuf[m.ID] = *m.Value
		logger.Logger.Debug(fmt.Sprintf("Gauge saved: %s=%s", m.ID, m.ValueStr()))

	default:
		return errors.New("wrong type")
	}
	return nil
}

func (repo *InMemMetricsRepository) GetAll() []models.Metrics {
	var metrics []models.Metrics

	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()
	for cName, cValue := range repo.counterBuf {
		value := float64(cValue)
		metrics = append(metrics, models.Metrics{ID: cName, MType: models.Counter, Value: &value})
	}
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()
	for gName, gValue := range repo.gaugeBuf {
		metrics = append(metrics, models.Metrics{ID: gName, MType: models.Gauge, Value: &gValue})
	}

	return metrics
}
