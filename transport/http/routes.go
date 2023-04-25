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

type HTTPHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) *Error
}

type HTTPHandlerFunc func(w http.ResponseWriter, r *http.Request) *Error

func (h HTTPHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) *Error {
	return h(w, r)
}

func Handler(h HTTPHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiErr := h.ServeHTTP(w, r); apiErr != nil {
			RespondJSON(w, apiErr, apiErr.Code)
		}
	})
}

type Route interface {
	// Handler returns the http handler.
	Handler() HTTPHandler
	// Method returns the http method that the route responds to.
	Method() string
	// Path returns the subpath where the route responds to.
	Path() string
}

func RespondJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
