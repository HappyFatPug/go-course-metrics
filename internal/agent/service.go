package agent

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
)

type AgentService struct {
	repo                 repository.MetricsRepository
	client               *Client
	memStatNamesToUpdate string
}

func NewAgentService(repo repository.MetricsRepository, client *Client) *AgentService {
	return &AgentService{
		repo:   repo,
		client: client,
		memStatNamesToUpdate: `Alloc BuckHashSys Frees GCCPUFraction GCSys HeapAlloc HeapIdle 
		HeapInuse HeapObjects HeapReleased HeapSys LastGC Lookups MCacheInuse MCacheSys 
		MSpanInuse MSpanSys Mallocs NextGC NumForcedGC NumGC OtherSys PauseTotalNs StackInuse 
		StackSys Sys TotalAlloc`,
	}
}

func (s *AgentService) UpdateAll(pollCount int64, rValue float64) error {
	if err := s.updateMemStatMetrics(); err != nil {
		return err
	}

	s.repo.Save(&models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &pollCount})
	s.repo.Save(&models.Metrics{ID: "RandomValue", MType: models.Gauge, Value: &rValue})

	return nil
}

func (s *AgentService) updateMemStatMetrics() error {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	val := reflect.ValueOf(memStat)

	for _, fieldName := range strings.Fields(s.memStatNamesToUpdate) {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("Failed to read mem stat field: %s", fieldName)
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
		s.repo.Save(&models.Metrics{ID: fieldName, MType: models.Gauge, Value: &value})
	}

	return nil
}

func (s *AgentService) SendAll() {
	err := s.SendMemStatMetrics()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	err = s.sendMetrics(models.Counter, "PollCount")
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	err = s.sendMetrics(models.Gauge, "RandomValue")
	if err != nil {
		logger.Logger.Error(err.Error())
	}

}

func (s *AgentService) SendMemStatMetrics() error {
	for _, mName := range strings.Fields(s.memStatNamesToUpdate) {
		err := s.sendMetrics(models.Gauge, mName)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}
	return nil
}

func (s *AgentService) sendMetrics(mType, name string) error {
	m, ok := s.readMetrics(mType, name)
	if !ok {
		logger.Logger.Warn(fmt.Sprintf("Cannot read metrics: %s", name))
		return nil
	}
	return s.client.postMetrics(m)
}

func (s *AgentService) readMetrics(mType string, name string) (*models.Metrics, bool) {
	return s.repo.Find(mType, name)
}
