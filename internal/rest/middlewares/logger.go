package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	log *zap.Logger
}

type responseData struct {
	status int
	size int
}

type loggerResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggerResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) 
	r.responseData.size += size
	return size, err
}

func (r *loggerResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) 
	r.responseData.status = statusCode
}

func NewLoggerMiddleware(log *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		log: log,
	}
}

func (lm *LoggerMiddleware) Logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		respData := &responseData{
			status: 0,
			size: 0,
		}
		lw := loggerResponseWriter{
			ResponseWriter: w,
			responseData: respData,
		}

		h.ServeHTTP(&lw, r)

		timeSince := time.Since(timeStart)

		lm.log.Info(
			"request info",
			zap.String("URI", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", timeSince),
		)
	}

	return http.HandlerFunc(fn)
}