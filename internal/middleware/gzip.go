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
	if c.useCompression {
		compressed, err := compress.Compress(p)
		if err != nil {
			return 0, err
		}
		_, err = c.writer.Write(compressed)
		return len(p), err
	}

	return c.writer.Write(p)
}

func (c *compressResponseWriter) WriteHeader(statusCode int) {
	contentType := c.rw.Header().Get("Content-Type")
	if slices.Contains(c.supportedContentTypes, strings.Split(contentType, ";")[0]) {
		c.rw.Header().Set("Content-Encoding", "gzip")
		c.useCompression = true
	}
	c.writer = c.rw

	c.headerSet = true
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
