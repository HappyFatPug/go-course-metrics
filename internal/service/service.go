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
		err = errors.New("wrong metric type")
	}

	return err
}

func (s *MetricsService) updateCounter(name, strValue string) error {
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}
	s.repo.Save(&models.Metrics{ID: name, MType: models.Counter, Delta: &value})
	logger.Logger.Debug(fmt.Sprintf("Counter value %s updated, new value: %d", name, value))
	// logger.Logger.Debug(fmt.Sprintf("Received counter value: %d", value))
	return nil
}

func (s *MetricsService) updateGauge(name string, strValue string) error {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return err
	}
	s.repo.Save(&models.Metrics{ID: name, MType: models.Gauge, Value: &value})
	logger.Logger.Debug(fmt.Sprintf("Gauge value %s updated, new value: %.2f", name, value))
	// logger.Logger.Debug(fmt.Sprintf("Received gauge value: %.2f", value))
	return nil
}

func (s *MetricsService) GetMetric(mType, mName string) (*models.Metrics, error) {
	metrics, ok := s.repo.Find(mType, mName)
	if !ok {
		return nil, errors.New("Metrics not found")
	}
	return metrics, nil
}
