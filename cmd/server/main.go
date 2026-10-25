package main

import (
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/HappyFatPug/go-course-metrics/internal/repository"
	"github.com/HappyFatPug/go-course-metrics/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	logger.ConfigureLogger()

	s := service.NewMetricsService(repository.NewInMemMetricsRepository())

	router := mux.NewRouter()
	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(s)).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
