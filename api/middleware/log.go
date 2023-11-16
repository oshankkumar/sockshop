package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/oshankkumar/sockshop/api/httpkit"
)

func WithLog(log *zap.Logger) httpkit.MiddlewareFunc {
	return func(method, pattern string, h httpkit.Handler) httpkit.Handler {
		return httpkit.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *httpkit.Error {
			start := time.Now()

			wr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			apiErr := h.ServeHTTP(wr, r)

			log.Info("request served",
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.Int("status", wr.Status()),
				zap.Int("bytes_written", wr.BytesWritten()),
				zap.Duration("took", time.Since(start)),
			)

			if apiErr != nil {
				log.Error("error in serving req", zap.Error(apiErr.Err))
			}

			return apiErr
		})
	}
}
