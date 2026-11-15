package app

import (
	"log/slog"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
	"github.com/HappyFatPug/go-course-metrics/internal/service"
	"github.com/HappyFatPug/go-course-metrics/internal/store"
)

type App struct {
	Debug   bool
	Service *service.MetricsService
	Logger  *slog.Logger
}

func ConfigureTestApp(debug bool, storeCfg *config.StoreConfig) *App {
	var lLevel slog.Leveler = slog.LevelInfo
	if debug {
		lLevel = slog.LevelDebug
	}
	logger := logger.ConfigureLogger(lLevel)

	memRepo := repository.NewInMemMetricsRepository()

	store := store.NewMetricsStore(storeCfg, logger)

	return &App{
		Debug:   debug,
		Service: service.NewMetricsService(memRepo, store, logger),
		Logger:  logger,
	}
}
