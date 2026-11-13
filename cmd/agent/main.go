package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/agent"
	"github.com/HappyFatPug/go-course-metrics/internal/config"
)

func main() {
	config := config.GetAgentConfig()
	app := agent.ConfigureTestApp(config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info("Poller started", slog.String("period", app.PollInterval.String()))
		worker(ctx, app.PollInterval, app.Service.UpdateAll)
		app.Logger.Info("Poller stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info("Sender started", slog.String("period", app.ReportInterval.String()))
		worker(ctx, app.ReportInterval, app.Service.SendAll)
		app.Logger.Info("Sender stopped")
	}()

	<-ctx.Done()
	app.Logger.Warn("Shutdown signal received...")
	wg.Wait()
	app.Logger.Warn("Service stopped")
}

func worker(ctx context.Context, interval time.Duration, fun func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fun()
		}
	}
}
