package main

import (
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	"github.com/HappyFatPug/go-course-metrics/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	app := app.ConfigureTestApp(false)

	router := mux.NewRouter()
	router.HandleFunc("/", handler.GetAll(app)).Methods("GET")
	router.HandleFunc("/update/{type}/{name}/{value}", handler.UpdateMetric(app)).Methods("POST")
	router.HandleFunc("/value/{type}/{name}", handler.GetMetric(app)).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		app.Logger.Error(err.Error())
	}
}
