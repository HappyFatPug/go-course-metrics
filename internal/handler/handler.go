package handler

import (
	"net/http"

	"github.com/HappyFatPug/go-course-metrics/internal/service"
	"github.com/gorilla/mux"
)

func UpdateMetric(s *service.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["type"]
		metricName := vars["name"]
		metricValue := vars["value"]

		// fmt.Printf("%s %s %s\n", metricName, metricType, metricValue)

		err := s.UpdateMetric(metricType, metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
