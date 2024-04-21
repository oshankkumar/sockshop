package httpkit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func DiscardErr(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = h.ServeHTTP(w, r)
	})
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e Error) Error() string {
	msg := fmt.Sprintf("Error %d: %s", e.Code, e.Message)
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
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

type MiddlewareFunc func(method, pattern string, h Handler) Handler

func ChainMiddleware(mm ...MiddlewareFunc) MiddlewareFunc {
	return func(method, pattern string, h Handler) Handler {
		for i := range mm {
			h = mm[len(mm)-i-1](method, pattern, h)
		}
		return h
	}
}
