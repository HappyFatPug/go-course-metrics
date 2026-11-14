package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

func WithLogging(logger *slog.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
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

		return http.HandlerFunc(logFn)
	}
}
