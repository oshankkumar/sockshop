package http

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) *Error
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) *Error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) *Error {
	return h(w, r)
}

func ToStdHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiErr := h.ServeHTTP(w, r); apiErr != nil {
			RespondJSON(w, apiErr, apiErr.Code)
		}
	})
}

type Router interface {
	Routes() []Route
}

type Route struct {
	Method  string
	Path    string
	Handler Handler
}

func RespondJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
