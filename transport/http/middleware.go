package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Middleware func(Handler) Handler

func WithLogging(log *zap.Logger) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
			if apiErr := h.ServeHTTP(w, r); apiErr != nil {
				log.Error("api error", zap.Any("details", apiErr), zap.Error(apiErr.Err))
				return apiErr
			}

			log.Info("api success", zap.String("url", r.RequestURI))
			return nil
		})
	}
}

func ChainMiddleware(mm ...Middleware) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
			for _, m := range mm {
				h = m(h)
			}
			return h.ServeHTTP(w, r)
		})
	}
}
