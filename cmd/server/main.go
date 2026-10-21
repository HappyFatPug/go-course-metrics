package main

import (
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	"github.com/gorilla/mux"
)

func main() {
	logger.ConfigureLogger()

	r := mux.NewRouter()
	r.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
