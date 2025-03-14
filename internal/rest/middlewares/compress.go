package middlewares

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type CompressMiddleware struct {
	log *zap.Logger
}

func NewCompressMiddleware(log *zap.Logger) *CompressMiddleware {
	return &CompressMiddleware{
		log: log,
	}
}

func (m *CompressMiddleware) Compress(h http.Handler) http.Handler {
	const op = "middlewares.Compress"

	fn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			log.Println("ГЗИП СТОИТ")
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
					m.log.Info("cannot create compress reader",
						zap.String("path", op),
						zap.Error(err),
					)

					w.WriteHeader(http.StatusInternalServerError)
					return
			}
			r.Body = cr

			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(fn)
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
			w:  w,
			zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
			c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
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