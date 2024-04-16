package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oshankkumar/sockshop/api/httpkit"
)

func NewMux() *HTTPMux {
	return &HTTPMux{router: chi.NewMux()}
}

type HTTPMux struct {
	router chi.Router
}

func (m *HTTPMux) Method(method, pattern string, h httpkit.Handler) {
	m.router.Method(method, pattern, httpkit.ToStdHandler(h))
}

func (m *HTTPMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func NewGoMux() *GoMux {
	return &GoMux{mux: http.NewServeMux()}
}

type GoMux struct {
	mux *http.ServeMux
}

func (m *GoMux) Method(method, pattern string, h httpkit.Handler) {
	m.mux.Handle(method+" "+pattern, httpkit.ToStdHandler(h))
}

func (m *GoMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
