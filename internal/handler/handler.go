package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/gorilla/mux"
)

func UpdateMetric(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["type"]
		metricName := vars["name"]
		metricValue := vars["value"]

		err := app.Service.UpdateMetric(metricType, metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func GetMetric(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["type"]
		metricName := vars["name"]

		m, ok := app.Service.GetMetric(metricType, metricName)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(m.ValueStr()))
	}
}

type allMetricsTeplateItem struct {
	Name, Value string
}

type allMetricsTemplateData struct {
	Counters, Gauges []allMetricsTeplateItem
}

var allMetricsT *template.Template

func init() {
	var err error
	allMetricsT, err = template.ParseFiles("../../web/templates/all-metrics.html")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetAll(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tData := allMetricsTemplateData{}
		for _, m := range app.Service.GetAllMetrics() {
			switch m.MType {
			case models.Counter:
				tData.Counters = append(tData.Counters, allMetricsTeplateItem{Name: m.ID, Value: m.ValueStr()})
			case models.Gauge:
				tData.Gauges = append(tData.Gauges, allMetricsTeplateItem{Name: m.ID, Value: m.ValueStr()})
			}
		}

		var body bytes.Buffer
		err := allMetricsT.Execute(&body, tData)
		if err != nil {
			logger.Logger.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		body.WriteTo(w)
	}
}
