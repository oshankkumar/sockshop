package httpkit

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func RespondError(w http.ResponseWriter, apiErr *Error) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(apiErr.Code)
	_ = json.NewEncoder(w).Encode(apiErr)
}

func RespondJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type MiddlewareFunc func(method, pattern string, h http.Handler) http.Handler

func ChainMiddleware(mm ...MiddlewareFunc) MiddlewareFunc {
	return func(method, pattern string, h http.Handler) http.Handler {
		for i := range mm {
			h = mm[len(mm)-i-1](method, pattern, h)
		}
		return h
	}
}
