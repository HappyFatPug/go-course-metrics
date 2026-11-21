package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/HappyFatPug/go-course-metrics/pkg/compress"
	"github.com/gorilla/mux"
)

type compressResponseWriter struct {
	rw                    http.ResponseWriter
	writer                io.Writer
	supportedContentTypes []string
	headerSet             bool
	useCompression        bool
}

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		rw:                    w,
		supportedContentTypes: []string{"application/json", "text/html"},
	}
}

func (c *compressResponseWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *compressResponseWriter) Write(p []byte) (int, error) {
	if !c.headerSet {
		c.WriteHeader(http.StatusOK)
	}

	return c.writer.Write(p)
}

func (c *compressResponseWriter) WriteHeader(statusCode int) {
	if c.headerSet {
		return
	}
	c.headerSet = true
	c.writer = c.rw

	contentType := strings.TrimSpace(
		strings.Split(c.Header().Get("Content-Type"), ";")[0],
	)

	if statusCode >= 200 &&
		statusCode < 300 &&
		statusCode != http.StatusNoContent &&
		statusCode != http.StatusResetContent &&
		c.Header().Get("Content-Encoding") == "" &&
		slices.Contains(c.supportedContentTypes, contentType) {

		c.Header().Set("Content-Encoding", "gzip")
		c.Header().Del("Content-Length")
		c.writer = compress.NewCompressWriter(c.rw)
	}
	c.rw.WriteHeader(statusCode)
}

func (c *compressResponseWriter) Close() error {
	if closer, ok := c.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func GzipMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w

			contentEncoding := r.Header.Get("Content-Encoding")
			if strings.Contains(contentEncoding, "gzip") {
				decBody, err := compress.Decompress(r.Body)
				defer r.Body.Close()
				if err != nil {
					logger.Error("gzip read error", slog.String("msg", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(decBody))
			}

			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				cw := newCompressResponseWriter(w)
				ow = cw
				defer cw.Close()
			}

			h.ServeHTTP(ow, r)
		})
	}
}
