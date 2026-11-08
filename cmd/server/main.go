package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	"github.com/HappyFatPug/go-course-metrics/internal/config"
	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	config := config.GetServerConfig()

	app := app.ConfigureTestApp(config.Debug)

	router := mux.NewRouter()
	router.HandleFunc("/", handler.GetAll(app.Service, app.Logger)).Methods("GET")
	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(app.Service, app.Logger)).Methods("POST")
	router.HandleFunc("/value/{type}/{name}", handler.GetMetric(app.Service)).Methods("GET")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := http.Server{
		Addr:    config.NetAddr.String(),
		Handler: router,
	}

	go func() {
		app.Logger.Info("Starting server", slog.String("URL", config.NetAddr.String()))
		err := server.ListenAndServe()
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}()

	<-ctx.Done()
	app.Logger.Warn("Shutdown signal received...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	} else {
		app.Logger.Warn("Server stopped")
	}
}
