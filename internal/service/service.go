package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
)

type MetricsService struct {
	repo repository.MetricsRepository
}

func NewMetricsService(repo repository.MetricsRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) UpdateMetric(mType, mName, mValue string) error {
	var err error
	switch mType {
	case models.Counter:
		err = s.updateCounter(mName, mValue)
	case models.Gauge:
		err = s.updateGauge(mName, mValue)
	default:
		err = errors.New("wrong type")
	}

	return err
}

func (s *MetricsService) updateCounter(name, strValue string) error {
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}
	metrics := s.GetOrCreateMetric(models.Counter, name)
	if metrics.Delta != nil {
		value += *metrics.Delta
	}
	metrics.Delta = &value

	s.repo.Save(metrics)
	logger.Logger.Debug(fmt.Sprintf("Counter value %s updated, new value: %d", name, value))
	return nil
}

func (s *MetricsService) updateGauge(name string, strValue string) error {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return err
	}

	s.repo.Save(&models.Metrics{ID: name, MType: models.Gauge, Value: &value})
	logger.Logger.Debug(fmt.Sprintf("Gauge value %s updated, new value: %.2f", name, value))
	return nil
}

func (s *MetricsService) GetOrCreateMetric(mType, name string) *models.Metrics {
	metrics, ok := s.repo.Find(mType, name)
	if !ok {
		return &models.Metrics{ID: name, MType: mType}
	}
	return metrics
}

func (s *MetricsService) GetMetric(mType, name string) (*models.Metrics, bool) {
	return s.repo.Find(mType, name)
}

func (s *MetricsService) GetAllMetrics() []models.Metrics {
	return s.repo.GetAll()
}
