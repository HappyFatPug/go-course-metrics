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

func (s *MetricsService) UpdateMetricRaw(mType, mName, mValue string) error {
	m, err := s.buildMetric(mType, mName, mValue)
	if err != nil {
		return err
	}

	return s.UpdateMetric(m)
}

func (s *MetricsService) buildMetric(mType, mName, mValue string) (*models.Metrics, error) {
	switch mType {
	case models.Counter:
		value, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return nil, errors.New("value must be integer")
		}
		return &models.Metrics{
			ID:    mName,
			MType: models.Counter,
			Delta: &value,
		}, nil
	case models.Gauge:
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return nil, err
		}
		return &models.Metrics{
			ID:    mName,
			MType: models.Gauge,
			Value: &value,
		}, nil
	default:
		return nil, errors.New("wrong type")
	}
}

func (s *MetricsService) UpdateMetric(m *models.Metrics) error {
	var err error
	switch m.MType {
	case models.Counter:
		err = s.updateCounter(m)
	case models.Gauge:
		err = s.updateGauge(m)
	default:
		err = errors.New("wrong type")
	}

	return err
}

func (s *MetricsService) updateCounter(updateData *models.Metrics) error {
	if updateData.MType != models.Counter {
		return errors.New("wrong type")
	}

	dbMetrics := s.getOrCreateMetric(models.Counter, updateData.ID)
	if dbMetrics.Delta != nil {
		*updateData.Delta += *dbMetrics.Delta
	}
	dbMetrics.Delta = updateData.Delta

	err := s.repo.Save(dbMetrics)
	if err != nil {
		return err
	}
	s.logger.Debug("Metric updated", slog.String("mType", models.Counter), slog.String("ID", dbMetrics.ID),
		slog.Int64("value", *dbMetrics.Delta))
	return nil
}

func (s *MetricsService) updateGauge(updateData *models.Metrics) error {
	err := s.repo.Save(updateData)
	if err != nil {
		return err
	}
	s.logger.Debug("Metric updated", slog.String("mType", models.Gauge), slog.String("ID", updateData.ID),
		slog.Float64("value", *updateData.Value))
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
