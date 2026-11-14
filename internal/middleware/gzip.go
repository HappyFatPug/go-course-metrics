package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
)

type compressWriter struct {
	rw                    http.ResponseWriter
	writer                io.Writer
	supportedContentTypes []string
	useCompressor         bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		rw:                    w,
		supportedContentTypes: []string{"application/json", "text/html"},
	}
}

func (c *compressWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.useCompressor {
		c.writer = gzip.NewWriter(c.rw)
	} else {
		c.writer = c.rw
	}

	return c.writer.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	contentType := c.rw.Header().Get("Content-Type")
	if slices.Contains(c.supportedContentTypes, strings.Split(contentType, ";")[0]) {
		c.rw.Header().Set("Content-Encoding", "gzip")
		c.useCompressor = true
	}
	c.rw.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if closer, ok := c.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w

			contentEncoding := r.Header.Get("Content-Encoding")
			if strings.Contains(contentEncoding, "gzip") {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					logger.Error("gzip read error", slog.String("msg", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}

			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			h.ServeHTTP(ow, r)
		})
	}
}
