package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/app"
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

		if app.Debug {
			mm, err := app.Service.GetMetric(metricType, metricName)
			if err != nil {
				app.Logger.Error(err.Error())
			} else {
				jsonData, err := json.Marshal(mm)
				if err != nil {
					app.Logger.Error(err.Error())
				} else {
					app.Logger.Debug(fmt.Sprintf("Updated metrics: %s", jsonData))
				}
			}
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
