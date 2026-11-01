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

func ConfigureTestApp(pollInterval, reportInterval int64, host string, port uint16, debug bool) *App {
	memRepo := repository.NewInMemMetricsRepository()

	client := NewClient("http", host, uint(port))

	var lLevel slog.Leveler = slog.LevelInfo
	if debug {
		lLevel = slog.LevelDebug
	}

	return &App{
		Debug:          debug,
		Service:        NewAgentService(memRepo, client),
		Logger:         logger.ConfigureLogger(lLevel),
		PollInterval:   time.Second * time.Duration(pollInterval),
		ReportInterval: time.Second * time.Duration(reportInterval),
	}
}
