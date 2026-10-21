package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func getLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler)
}

func ConfigureLogger() {
	Logger = getLogger()
}
