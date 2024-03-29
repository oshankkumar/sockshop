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

func RespondJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type Router interface {
	InstallRoutes(mux Mux)
}

type RouterFunc func(Mux)

func (r RouterFunc) InstallRoutes(mux Mux) { r(mux) }

type Mux interface {
	http.Handler
	Method(method, pattern string, h http.Handler)
}

type Routers []Router

func (rr Routers) InstallRoutes(mux Mux) {
	for _, r := range rr {
		r.InstallRoutes(mux)
	}
}

type InstrumentedMux struct {
	Mux
	middleware MiddlewareFunc
}

func (i *InstrumentedMux) Method(method, pattern string, h http.Handler) {
	i.Mux.Method(method, pattern, i.middleware(method, pattern, h))
}

func NewInstrumentedMux(m Mux, middleware MiddlewareFunc) Mux {
	return &InstrumentedMux{Mux: m, middleware: middleware}
}
