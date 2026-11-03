package handler

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/gorilla/mux"
)

type Service interface {
	UpdateMetric(mType, mName, mValue string) error
	GetMetric(mType, name string) (*models.Metrics, bool)
	GetAllMetrics() []models.Metrics
}

func UpdateMetric(srv Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["type"]
		metricName := vars["name"]
		metricValue := vars["value"]

		if metricType != models.Counter && metricType != models.Gauge {
			http.Error(w, "Metric type should be counter/gauge", http.StatusBadRequest)
			return
		}

		err := srv.UpdateMetric(metricType, metricName, metricValue)
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func GetMetric(srv Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		metricType := vars["type"]
		metricName := vars["name"]

		m, ok := srv.GetMetric(metricType, metricName)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(m.ValueStr()))
	}
}

func GetAll(srv Service, logger *slog.Logger) http.HandlerFunc {
	tmpl := template.Must(template.New("getAllTemplate").Parse(`<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
			</head>
			<body>
			{{range .Counters}}
				<div>
					<span>{{.Name}}</span>: <span>{{.Value}}</span>
				</div>
			{{end}}
			{{range .Gauges}}
				<div>
					<span>{{.Name}}</span>: <span>{{.Value}}</span>
				</div>
			{{end}}   
			</body>
		</html>`))

	return func(w http.ResponseWriter, r *http.Request) {
		type allMetricsTeplateItem struct {
			Name, Value string
		}

		type allMetricsTemplateData struct {
			Counters, Gauges []allMetricsTeplateItem
		}

		tData := allMetricsTemplateData{}
		for _, m := range srv.GetAllMetrics() {
			switch m.MType {
			case models.Counter:
				tData.Counters = append(tData.Counters, allMetricsTeplateItem{Name: m.ID, Value: m.ValueStr()})
			case models.Gauge:
				tData.Gauges = append(tData.Gauges, allMetricsTeplateItem{Name: m.ID, Value: m.ValueStr()})
			}
		}

		var body bytes.Buffer
		err := tmpl.Execute(&body, tData)
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		body.WriteTo(w)
	}
}
