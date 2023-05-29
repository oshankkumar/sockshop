package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oshankkumar/sockshop/api"
)

type APIServer struct {
	ImagePath     string
	HealthChecker HealthChecker
	Middleware    Middleware
}

func (s *APIServer) CreateMux(routers ...Router) *mux.Router {
	serveMux := mux.NewRouter()

	for _, router := range routers {
		for _, route := range router.Routes() {
			h := s.Middleware(route.Handler)
			serveMux.Methods(route.Method).Path(route.Path).Handler(ToStdHandler(h))
		}
	}

	imgServer := http.FileServer(http.Dir(s.ImagePath))
	imgServer = http.StripPrefix("/catalogue/images/", imgServer)

	serveMux.Methods(http.MethodGet).PathPrefix("/catalogue/images/").Handler(imgServer)
	serveMux.Methods(http.MethodGet).Path("/health").Handler(ToStdHandler(HealthCheckHandler(s.HealthChecker)))

	return serveMux
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
