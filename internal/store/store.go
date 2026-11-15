package store

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type MetricsStore struct {
	Sync     bool
	filename string
	logger   *slog.Logger
	mu       sync.Mutex
}

func NewMetricsStore(cfg *config.StoreConfig, logger *slog.Logger) *MetricsStore {
	return &MetricsStore{Sync: cfg.StoreInterval == 0, filename: cfg.FileStoragePath, logger: logger}
}

func (s *MetricsStore) IsSync() bool {
	return s.Sync
}

func (s *MetricsStore) Dump(metrics []models.Metrics) error {
	if s.filename == "" {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	for _, m := range metrics {
		err := enc.Encode(&m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricsStore) ListStored() ([]models.Metrics, error) {
	if s.filename == "" {
		return nil, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metrics []models.Metrics
	dec := json.NewDecoder(file)

	for {
		var m models.Metrics
		err := dec.Decode(&m)
		if err != nil {
			s.logger.Debug("store.ListStored err:", slog.String("err", err.Error()))
			break
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
