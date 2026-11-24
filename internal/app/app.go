package app

import (
	"log/slog"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	"github.com/HappyFatPug/go-course-metrics/internal/db"
	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
	"github.com/HappyFatPug/go-course-metrics/internal/service"
	"github.com/HappyFatPug/go-course-metrics/internal/store"
)

type App struct {
	Debug      bool
	Service    *service.MetricsService
	Repository *repository.StorableMetricsRepository
	DB         *db.DB
	Logger     *slog.Logger
}

func ConfigureApp(debug bool, storeCfg *config.StoreConfig, dbCfg *config.DbConfig) (*App, error) {
	var lLevel slog.Leveler = slog.LevelInfo
	if debug {
		lLevel = slog.LevelDebug
	}
	logger := logger.ConfigureLogger(lLevel)

	memRepo := repository.NewInMemMetricsRepository()
	store := store.NewMetricsStore(storeCfg.FileStoragePath, logger)
	rep, err := repository.NewStorableMetricsRepository(
		memRepo,
		store,
		storeCfg,
	)
	if err != nil {
		return nil, err
	}

	db, err := db.NewDB(dbCfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &App{
		Debug:      debug,
		Service:    service.NewMetricsService(rep, logger),
		Repository: rep,
		DB:         db,
		Logger:     logger,
	}, nil
}
