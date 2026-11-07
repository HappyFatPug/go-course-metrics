package agent

import (
	"log/slog"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
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

func ConfigureTestApp(config *config.AgentConfig) *App {
	var lLevel slog.Leveler = slog.LevelInfo
	if config.Debug {
		lLevel = slog.LevelDebug
	}
	logger := logger.ConfigureLogger(lLevel)

	memRepo := repository.NewInMemMetricsRepository()

	client := NewClient(config.ClientConfig, logger)

	return &App{
		Debug:          config.Debug,
		Service:        NewAgentService(memRepo, client, logger),
		Logger:         logger,
		PollInterval:   time.Second * time.Duration(config.PollInterval),
		ReportInterval: time.Second * time.Duration(config.ReportInterval),
	}
}
