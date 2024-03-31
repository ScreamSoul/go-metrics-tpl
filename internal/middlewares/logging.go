package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (lm *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lm.logger.Info("Request received",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
		)

		lw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lw, r)

		lm.logger.Info("Response sent",
			zap.Int("status", lw.statusCode),
			zap.Int("size", lw.size),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(b)
	lw.size += size
	return size, err
}
