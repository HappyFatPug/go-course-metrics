package logger

import (
	"log/slog"
	"os"
)

func ConfigureLogger(level slog.Leveler) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
