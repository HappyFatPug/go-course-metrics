package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	parseFlags()

	app := app.ConfigureTestApp(config.debug)

	router := mux.NewRouter()
	router.HandleFunc("/", handler.GetAll(app)).Methods("GET")
	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(app)).Methods("POST")
	router.HandleFunc("/value/{type}/{name}", handler.GetMetric(app)).Methods("GET")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := http.Server{
		Addr:    config.netAddr.String(),
		Handler: router,
	}

	go func() {
		app.Logger.Info(fmt.Sprintf("Listennig %s", config.netAddr.String()))
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
		app.Logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	} else {
		app.Logger.Warn("Server stopped")
	}
}
