package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type MiddlewareFunc func(method, pattern string, h http.Handler) http.Handler

func LogginMiddleware(log *zap.Logger) MiddlewareFunc {
	return func(method, pattern string, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			h.ServeHTTP(wr, r)

			log.Info("request served",
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.Int("status", wr.Status()),
				zap.Int("bytes_written", wr.BytesWritten()),
				zap.Duration("took", time.Since(start)),
			)
		})
	}
}
