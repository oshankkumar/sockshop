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
