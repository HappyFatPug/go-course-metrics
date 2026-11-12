package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	appErrors "github.com/HappyFatPug/go-course-metrics/internal/errors"
)

func WithLogging(h http.Handler, logger *slog.Logger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Info(
			"http.Request",
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.Int("status", responseData.status),
			slog.Duration("duration", duration),
			slog.Int("size", responseData.size),
		)
	}

	return logFn
}

func JSONHandler(h http.Handler, logger *slog.Logger) http.HandlerFunc {
	checkJSONFn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") != "application/json" {
			resp := appErrors.JSONError{
				Code:    http.StatusUnsupportedMediaType,
				Message: http.StatusText(http.StatusUnsupportedMediaType),
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				logger.Debug(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		h.ServeHTTP(w, r)
	}

	return checkJSONFn
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
