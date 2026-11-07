package agent

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"reflect"
	"runtime"
	"strings"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type MetricsRepository interface {
	Find(mType, name string) (*models.Metrics, bool)
	Save(m *models.Metrics) error
	GetAll() []models.Metrics
}

type AgentService struct {
	repo                 MetricsRepository
	client               *Client
	memStatNamesToUpdate string
	logger               *slog.Logger
}

func NewAgentService(repo MetricsRepository, client *Client, logger *slog.Logger) *AgentService {
	return &AgentService{
		repo:   repo,
		client: client,
		memStatNamesToUpdate: `Alloc BuckHashSys Frees GCCPUFraction GCSys HeapAlloc HeapIdle 
		HeapInuse HeapObjects HeapReleased HeapSys LastGC Lookups MCacheInuse MCacheSys 
		MSpanInuse MSpanSys Mallocs NextGC NumForcedGC NumGC OtherSys PauseTotalNs StackInuse 
		StackSys Sys TotalAlloc`,
		logger: logger,
	}
}

func (s *AgentService) saveCounter(name string, value int64) error {
	fvalue := float64(value)
	return s.repo.Save(&models.Metrics{ID: name, MType: models.Counter, Value: &fvalue})
}

func (s *AgentService) saveGauge(name string, value float64) error {
	return s.repo.Save(&models.Metrics{ID: name, MType: models.Gauge, Value: &value})
}

func (s *AgentService) UpdateAll() {
	if err := s.updateMemStatMetrics(); err != nil {
		s.logger.Error(err.Error())
	}

	err := s.incCounter("PollCount")
	if err != nil {
		s.logger.Error(err.Error())
	}
	err = s.saveGauge("RandomValue", rand.Float64())
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *AgentService) incCounter(name string) error {
	metric, ok := s.readMetrics(models.Counter, name)
	if !ok {
		val := float64(0)
		metric = &models.Metrics{
			MType: models.Counter,
			ID:    name,
			Value: &val,
		}
	}

	value := *metric.Value + 1
	metric.Value = &value

	return s.saveMetrics(metric)
}

func (s *AgentService) updateMemStatMetrics() error {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	val := reflect.ValueOf(memStat)

	for _, fieldName := range strings.Fields(s.memStatNamesToUpdate) {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("failed to read mem stat field: %s", fieldName)
		}

		var value float64

		switch field.Kind() {
		case reflect.Float64, reflect.Float32:
			value = field.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = float64(field.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = float64(field.Int())
		}
		err := s.saveGauge(fieldName, value)
		if err != nil {
			return err
		}
		s.logger.Debug("MemStat gauge saved", slog.String("ID", fieldName), slog.Float64("value", value))
	}

	return nil
}

func (s *AgentService) SendAll() {
	err := s.SendMemStatMetrics()
	if err != nil {
		s.logger.Error(err.Error())
	}

	err = s.sendMetrics(models.Gauge, "RandomValue")
	if err != nil {
		s.logger.Error(err.Error())
	}

	err = s.sendMetrics(models.Counter, "PollCount")
	if err != nil {
		s.logger.Error(err.Error())
	} else {
		s.saveCounter("PollCount", 0)
	}
}

func (s *AgentService) SendMemStatMetrics() error {
	for _, mName := range strings.Fields(s.memStatNamesToUpdate) {
		err := s.sendMetrics(models.Gauge, mName)
		if err != nil {
			s.logger.Error(err.Error())
		}
	}
	return nil
}

func (s *AgentService) sendMetrics(mType, name string) error {
	m, ok := s.readMetrics(mType, name)
	if !ok {
		s.logger.Warn("Cannot read metrics", slog.String("ID", name))
		return nil
	}
	return s.client.postMetrics(m)
}

func (s *AgentService) readMetrics(mType string, name string) (*models.Metrics, bool) {
	return s.repo.Find(mType, name)
}

func (s *AgentService) saveMetrics(m *models.Metrics) error {
	return s.repo.Save(m)
}
