package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	appErrors "github.com/HappyFatPug/go-course-metrics/internal/errors"
)

func JSONRequired(h http.Handler, logger *slog.Logger) http.HandlerFunc {
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
