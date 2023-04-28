package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Middleware func(HTTPHandler) HTTPHandler

func WithLogging(log *zap.Logger) Middleware {
	return func(h HTTPHandler) HTTPHandler {
		return HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
			apiErr := h.ServeHTTP(w, r)
			if apiErr != nil {
				log.Error("api error", zap.Any("details", apiErr))
			}
			return apiErr
		})
	}
}

func ChainMiddleware(mm ...Middleware) Middleware {
	return func(h HTTPHandler) HTTPHandler {
		return HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) *Error {
			for _, m := range mm {
				h = m(h)
			}
			return h.ServeHTTP(w, r)
		})
	}
}
