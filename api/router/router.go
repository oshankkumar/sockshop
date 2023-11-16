package router

import (
	"net/http"

	"github.com/oshankkumar/sockshop/api/httpkit"
)

type Mux interface {
	http.Handler
	Method(method, pattern string, h httpkit.Handler)
}

type InstrumentedMux struct {
	base       Mux
	middleware httpkit.MiddlewareFunc
}

func (i *InstrumentedMux) Method(method, pattern string, h httpkit.Handler) {
	i.base.Method(method, pattern, i.middleware(method, pattern, h))
}

func (i *InstrumentedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.base.ServeHTTP(w, r)
}

func NewInstrumentedMux(m Mux, middleware httpkit.MiddlewareFunc) Mux {
	return &InstrumentedMux{base: m, middleware: middleware}
}

type Router interface {
	InstallRoutes(mux Mux)
}

type RouterFunc func(Mux)

func (r RouterFunc) InstallRoutes(mux Mux) { r(mux) }

type Routers []Router

func (rr Routers) InstallRoutes(mux Mux) {
	for _, r := range rr {
		r.InstallRoutes(mux)
	}
}
