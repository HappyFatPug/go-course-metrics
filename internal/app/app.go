package app

import (
	"log/slog"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
	"github.com/HappyFatPug/go-course-metrics/internal/service"
)

type App struct {
	Debug   bool
	Service *service.MetricsService
	Logger  *slog.Logger
}

func ConfigureTestApp(debug bool) *App {
	memRepo := repository.NewInMemMetricsRepository()
	var lLevel slog.Leveler = slog.LevelInfo
	if debug {
		lLevel = slog.LevelDebug
	}

	return &App{
		Debug:   debug,
		Service: service.NewMetricsService(memRepo),
		Logger:  logger.ConfigureLogger(lLevel),
	}
}
