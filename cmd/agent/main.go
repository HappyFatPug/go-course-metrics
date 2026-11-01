package main

import (
	"context"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/agent"
)

var PollCount int64 = 0
var RandomValue float64

func main() {
	parseFlags()

	app := agent.ConfigureTestApp(
		config.pollInterval,
		config.reportInterval,
		config.netAddr.Host,
		config.netAddr.Port,
		config.debug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info("Poller started")
		worker(ctx, app.PollInterval, func() {
			PollCount += 1
			RandomValue = rand.Float64()
			app.Service.UpdateAll(PollCount, RandomValue)
		})
		app.Logger.Info("Poller stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info("Sender started")
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
