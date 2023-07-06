package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oshankkumar/sockshop/api"
	"go.uber.org/zap"
)

type APIServer struct {
	Mux           chi.Router
	ImagePath     string
	HealthChecker HealthChecker
	Middleware    Middleware
	Logger        *zap.Logger
}

func (a *APIServer) Start(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr, Handler: a}
	errc := make(chan error, 1)

	a.Logger.Info("starting app http server :9090")

	go func(errc chan<- error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		a.Logger.Info("received context cancellation; shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("http server shutdonwn: %w", err)
		}
		return nil
	}

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return shutdown(time.Second * 5)
	}
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
