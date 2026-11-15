package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	"github.com/HappyFatPug/go-course-metrics/internal/config"
	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/HappyFatPug/go-course-metrics/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	config := config.GetServerConfig()

	app := app.ConfigureTestApp(config.Debug, config.StoreCfg)

	router := mux.NewRouter()

	router.HandleFunc("/", handler.GetAll(app.Service, app.Logger)).Methods("GET")

	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(app.Service, app.Logger)).Methods("POST")
	router.HandleFunc(
		"/update/",
		middleware.JSONRequired(handler.UpdateMetricJSON(app.Service, app.Logger), app.Logger),
	).Methods("POST")

	router.HandleFunc("/value/{type}/{name}", handler.GetMetric(app.Service)).Methods("GET")
	router.HandleFunc(
		"/value/",
		middleware.JSONRequired(handler.GetMetricJSON(app.Service, app.Logger), app.Logger),
	).Methods("POST")

	server := http.Server{
		Addr: config.NetAddr.String(),
		Handler: middleware.GzipMiddleware(app.Logger)(
			middleware.WithLogging(app.Logger)(
				router,
			),
		),
	}

	if config.StoreCfg.Restore {
		err := app.Service.RestoreTmpFile()
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info("Starting server", slog.String("URL", config.NetAddr.String()))
		err := server.ListenAndServe()
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}()

	if config.StoreCfg.StoreInterval > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tickerPeriod := time.Duration(config.StoreCfg.StoreInterval) * time.Second
			ticker := time.NewTicker(tickerPeriod)
			defer ticker.Stop()

			app.Logger.Info("Dumper starting", slog.String("period", tickerPeriod.String()))
			for {
				select {
				case <-ctx.Done():
					app.Logger.Info("Dumper: shutdown signal received")
					app.Service.DumpTmpFile()
					return
				case <-ticker.C:
					app.Service.DumpTmpFile()
				}
			}
		}()
	}

	<-ctx.Done()
	app.Logger.Warn("Shutdown signal received...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	} else {
		app.Logger.Warn("Server stopped")
	}

	wg.Wait()
}
