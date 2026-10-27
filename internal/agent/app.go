package agent

import (
	"log/slog"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
)

type App struct {
	Debug          bool
	Service        *AgentService
	Logger         *slog.Logger
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func ConfigureTestApp(pollInterval, reportInterval time.Duration, debug bool) *App {
	memRepo := repository.NewInMemMetricsRepository()
	client := NewClient("http", "localhost", 8080)
	var lLevel slog.Leveler = slog.LevelInfo
	if debug {
		lLevel = slog.LevelDebug
	}

	return &App{
		Debug:          debug,
		Service:        NewAgentService(memRepo, client),
		Logger:         logger.ConfigureLogger(lLevel),
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}
}
