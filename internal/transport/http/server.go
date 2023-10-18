package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oshankkumar/sockshop/api"
	"go.uber.org/zap"
)

type APIServer struct {
	Addr   string
	Mux    http.Handler
	Logger *zap.Logger

	server *http.Server
}

func (a *APIServer) Start(ctx context.Context) error {
	a.server = &http.Server{Addr: a.Addr, Handler: a.Mux}
	errc := make(chan error, 1)

	a.Logger.Info("starting app http server :9090")

	go func(errc chan<- error) {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		a.Logger.Info("received context cancellation; shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := a.server.Shutdown(shutCtx); err != nil {
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

func HealthCheckRouter(hc HealthChecker) Router {
	return RouterFunc(func(m Mux) {
		m.Method(http.MethodGet, "/health", ToStdHandler(HealthCheckHandler(hc)))
	})
}

func ImageServeRouter(path string) Router {
	return RouterFunc(func(mux Mux) {
		mux.Method(http.MethodGet, "/catalogue/images/*", http.StripPrefix(
			"/catalogue/images/", http.FileServer(http.Dir(path)),
		))
	})
}
