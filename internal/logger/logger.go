package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func ConfigureLogger(level slog.Leveler) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	Logger = slog.New(handler)
	return Logger
}
