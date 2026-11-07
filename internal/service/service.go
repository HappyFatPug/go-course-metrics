package service

import (
	"errors"
	"log/slog"
	"strconv"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type MetricsRepository interface {
	Find(mType, name string) (*models.Metrics, bool)
	Save(m *models.Metrics) error
	GetAll() []models.Metrics
}

type MetricsService struct {
	repo   MetricsRepository
	logger *slog.Logger
}

func NewMetricsService(repo MetricsRepository, logger *slog.Logger) *MetricsService {
	return &MetricsService{
		repo:   repo,
		logger: logger,
	}
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
		return errors.New("value must be integer")
	}
	metrics := s.getOrCreateMetric(models.Counter, name)
	fValue := float64(value)
	if metrics.Value != nil {
		fValue = fValue + *metrics.Value
	}
	metrics.Value = &fValue

	err = s.repo.Save(metrics)
	if err != nil {
		return err
	}
	s.logger.Debug("Metric updated", slog.String("mType", models.Counter), slog.String("ID", name),
		slog.Int64("calue", value))
	return nil
}

func (s *MetricsService) updateGauge(name string, strValue string) error {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return err
	}

	err = s.repo.Save(&models.Metrics{ID: name, MType: models.Gauge, Value: &value})
	if err != nil {
		return err
	}
	s.logger.Debug("Metric updated", slog.String("mType", models.Gauge), slog.String("ID", name),
		slog.Float64("value", value))
	return nil
}

func (s *MetricsService) getOrCreateMetric(mType, name string) *models.Metrics {
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
