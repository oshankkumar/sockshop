package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oshankkumar/sockshop/api"
)

type APIServer struct {
	Mux           chi.Router
	ImagePath     string
	HealthChecker HealthChecker
	Middleware    Middleware
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *APIServer) InstallRoutes(routes ...Routes) {
	s.Mux.Get("/health", ToStdHandler(HealthCheckHandler(s.HealthChecker)).ServeHTTP)

	s.Mux.Handle("/catalogue/images/*", http.StripPrefix(
		"/catalogue/images/", http.FileServer(http.Dir(s.ImagePath)),
	))

	for _, route := range routes {
		for _, r := range route.Routes() {
			h := s.Middleware(r.Handler)
			s.Mux.Method(r.Method, r.Path, ToStdHandler(h))
		}
	}
}

type HealthChecker interface {
	CheckHealth(ctx context.Context) ([]api.Health, error)
}

type HealthCheckerFunc func(ctx context.Context) ([]api.Health, error)

func (h HealthCheckerFunc) CheckHealth(ctx context.Context) ([]api.Health, error) { return h(ctx) }

func HealthCheckHandler(h HealthChecker) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) *Error {
		hh, err := h.CheckHealth(r.Context())
		if err != nil {
			return &Error{http.StatusInternalServerError, err.Error(), err}
		}

		RespondJSON(w, api.HealthResponse{Healths: hh}, http.StatusOK)
		return nil
	}
}
