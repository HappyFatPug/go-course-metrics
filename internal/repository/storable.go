package repository

import (
	"log"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type Store interface {
	Dump([]models.Metrics) error
	ListStored() ([]models.Metrics, error)
}

type Repository interface {
	Find(mType, name string) (*models.Metrics, bool)
	Save(m *models.Metrics) error
	GetAll() []models.Metrics
}

type StorableMetricsRepository struct {
	Repository
	Store
	Sync bool
}

func NewStorableMetricsRepository(repo Repository, store Store, cfg *config.StoreConfig) *StorableMetricsRepository {
	r := StorableMetricsRepository{
		repo,
		store,
		cfg.StoreInterval == 0,
	}

	if cfg.Restore {
		err := r.RestoreDump()
		if err != nil {
			log.Fatal(err)
		}
	}

	return &r
}

func (r *StorableMetricsRepository) Save(m *models.Metrics) error {
	err := r.Repository.Save(m)
	if err != nil {
		return err
	}

	if r.Sync {
		return r.CreateDump()
	}

	return nil
}

func (r *StorableMetricsRepository) CreateDump() error {
	err := r.Dump(r.GetAll())
	if err != nil {
		return err
	}

	return nil
}

func (r *StorableMetricsRepository) RestoreDump() error {
	ms, err := r.ListStored()
	if err != nil {
		return err
	}
	for _, m := range ms {
		err := r.Repository.Save(&m)
		if err != nil {
			return err
		}
	}

	return nil
}
