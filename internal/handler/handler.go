package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	appErrors "github.com/HappyFatPug/go-course-metrics/internal/errors"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
	"github.com/gorilla/mux"
)

type Service interface {
	UpdateMetricRaw(mType, mName, mValue string) error
	UpdateMetric(*models.Metrics) error
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

		err := srv.UpdateMetricRaw(metricType, metricName, metricValue)
		if err != nil {
			logger.Debug(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateMetricJSON(srv Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			logger.Debug("JSON decode error", slog.String("msg", err.Error()))
			JSONErrorRequestHelper(w, r, logger, http.StatusBadRequest)
			return
		}

		if (m.ID == "" || m.MType == "" || m.Value == nil && m.Delta == nil) ||
			m.MType == models.Counter && m.Delta == nil ||
			m.MType == models.Gauge && m.Value == nil {
			JSONErrorRequestHelper(w, r, logger, http.StatusBadRequest)
			return
		}

		err := srv.UpdateMetric(&m)
		if err != nil {
			logger.Debug("UpdateMetric error", slog.String("msg", err.Error()))
			JSONErrorRequestHelper(w, r, logger, http.StatusBadRequest)
			return
		}

		var body bytes.Buffer
		if err := json.NewEncoder(&body).Encode(&m); err != nil {
			logger.Error("JSON encode error", slog.String("msg", err.Error()))
			JSONErrorRequestHelper(w, r, logger, http.StatusInternalServerError)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		body.WriteTo(w)
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

func GetMetricJSON(srv Service, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestM models.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&requestM); err != nil {
			logger.Debug("JSON decode error", slog.String("msg", err.Error()))
			JSONErrorRequestHelper(w, r, logger, http.StatusBadRequest)
			return
		}

		if requestM.ID == "" || requestM.MType == "" {
			JSONErrorRequestHelper(w, r, logger, http.StatusBadRequest)
			return
		}

		m, ok := srv.GetMetric(requestM.MType, requestM.ID)
		if !ok {
			JSONErrorRequestHelper(w, r, logger, http.StatusNotFound)
			return
		}

		var body bytes.Buffer
		if err := json.NewEncoder(&body).Encode(&m); err != nil {
			logger.Error("JSON encode error", slog.String("msg", err.Error()))
			JSONErrorRequestHelper(w, r, logger, http.StatusInternalServerError)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		body.WriteTo(w)
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
		w.WriteHeader(http.StatusOK)
		body.WriteTo(w)
	}
}

func JSONErrorRequestHelper(w http.ResponseWriter, r *http.Request, logger *slog.Logger, code int) {
	resp := appErrors.JSONError{
		Code:    code,
		Message: http.StatusText(code),
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(resp)
	if err != nil {
		logger.Error("JSON encode error", slog.String("msg", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	body.WriteTo(w)
}
