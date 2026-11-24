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

	app, err := app.ConfigureApp(config.Debug, config.StoreCfg, config.DBCfg)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", handler.GetAll(app.Service, app.Logger)).Methods("GET")
	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(app.Service, app.Logger)).Methods("POST")
	router.HandleFunc("/value/{type}/{name}", handler.GetMetric(app.Service)).Methods("GET")
	router.HandleFunc("/ping", handler.DBPing(app.DB, app.Logger)).Methods("GET")

	jsonRouter := router.PathPrefix("").Subrouter()
	jsonRouter.Use(middleware.JSONRequired(app.Logger))
	jsonRouter.HandleFunc("/value/", handler.GetMetricJSON(app.Service, app.Logger)).Methods("POST")
	jsonRouter.HandleFunc("/update/", handler.UpdateMetricJSON(app.Service, app.Logger)).Methods("POST")

	server := http.Server{
		Addr: config.NetAddr.String(),
		Handler: middleware.GzipMiddleware(app.Logger)(
			middleware.WithLogging(app.Logger)(
				router,
			),
		),
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
					err := app.Repository.CreateDump()
					if err != nil {
						app.Logger.Error("dump file create error", slog.String("error", err.Error()))
					}
					return
				case <-ticker.C:
					err := app.Repository.CreateDump()
					if err != nil {
						app.Logger.Error("dump file create error", slog.String("error", err.Error()))
					}
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

	dbShutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.DB.Close(dbShutdownCtx); err != nil {
		app.Logger.Error("DB close error", slog.String("error", err.Error()))
	} else {
		app.Logger.Warn("DB connection closed")
	}

	wg.Wait()
}
