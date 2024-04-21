package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/oshankkumar/sockshop/api/httpkit"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func WithLog(log *zap.Logger) httpkit.MiddlewareFunc {
	return func(method, pattern string, h httpkit.Handler) httpkit.Handler {
		return httpkit.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()

			wr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			err := h.ServeHTTP(wr, r)

			fields := []zap.Field{
				zap.String("method", method),
				zap.String("pattern", pattern),
				zap.String("url", r.RequestURI),
				zap.Int("status", wr.Status()),
				zap.Int("bytes_written", wr.BytesWritten()),
				zap.Duration("took", time.Since(start)),
			}
			if err != nil {
				log.Error("request failed", append(fields, zap.Error(err))...)
				return err
			}

			log.Info("request succeeded", fields...)
			return err
		})
	}
}

func WithHTTPErrStatus(method, pattern string, h httpkit.Handler) httpkit.Handler {
	return httpkit.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		err := h.ServeHTTP(w, r)
		if err == nil {
			return nil
		}

		var apiErr *httpkit.Error

		if !errors.As(err, &apiErr) {
			apiErr = &httpkit.Error{Code: http.StatusInternalServerError, Message: "something went wrong", Err: err}
		}

		httpkit.RespondError(w, apiErr)

		return err
	})
}
