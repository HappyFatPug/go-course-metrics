package store

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type MetricsStore struct {
	dir      string
	filename string
	logger   *slog.Logger
	mu       sync.Mutex
}

func NewMetricsStore(storagePath string, logger *slog.Logger) *MetricsStore {
	dir, filename := filepath.Split(storagePath)
	return &MetricsStore{
		filename: filename,
		dir:      dir,
		logger:   logger,
	}
}

func (s *MetricsStore) Dump(metrics []models.Metrics) error {
	if s.filename == "" {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tmpFile, err := os.CreateTemp(s.dir, ".tmp-*.json")
	if err != nil {
		return err
	}
	defer func() {
		if tmpFile != nil {
			if err := tmpFile.Close(); err != nil {
				s.logger.Error("close tmp file error", slog.String("err", err.Error()))
				return
			}
			if err := os.Remove(tmpFile.Name()); err != nil {
				s.logger.Error("remove tmp file error", slog.String("err", err.Error()))
				return
			}
		}
	}()

	enc := json.NewEncoder(tmpFile)

	for _, m := range metrics {
		err := enc.Encode(&m)
		if err != nil {
			return err
		}
	}

	if err := tmpFile.Close(); err != nil {
		s.logger.Error("close tmp file error", slog.String("err", err.Error()))
		return err
	}
	tmpFileName := tmpFile.Name()
	tmpFile = nil
	if err := os.Chmod(tmpFileName, 0666); err != nil {
		s.logger.Error("chmod tmp file error", slog.String("err", err.Error()))
		if removeErr := os.Remove(tmpFileName); removeErr != nil {
			s.logger.Error("remove tmp file error", slog.String("err", removeErr.Error()))
			return errors.Join(err, removeErr)
		}
		return err
	}
	if err := os.Rename(tmpFileName, filepath.Join(s.dir, s.filename)); err != nil {
		s.logger.Error("rename tmp file error", slog.String("err", err.Error()))
		if removeErr := os.Remove(tmpFileName); removeErr != nil {
			s.logger.Error("remove tmp file error", slog.String("err", removeErr.Error()))
			return errors.Join(err, removeErr)
		}
		return err
	}
	return nil
}

func (s *MetricsStore) ListStored() ([]models.Metrics, error) {
	if s.filename == "" {
		return nil, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(filepath.Join(s.dir, s.filename), os.O_RDONLY|os.O_CREATE, 0666)
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
