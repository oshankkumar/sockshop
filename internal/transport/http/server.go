package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/internal/domain"
	"go.uber.org/zap"
)

type APIServer struct {
	mux           chi.Router
	ImagePath     string
	HealthChecker HealthChecker
	Logger        *zap.Logger
	UserService   api.UserService
	SockLister    SockLister
	SockStore     domain.SockStore
}

func (s *APIServer) Start(ctx context.Context, addr string) error {
	s.mux = s.initRoutes()

	srv := &http.Server{Addr: addr, Handler: s}
	errc := make(chan error, 1)

	s.Logger.Info("starting app http server :9090")

	go func(errc chan<- error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		s.Logger.Info("received context cancellation; shutting down server")

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
	s.mux.ServeHTTP(w, r)
}

func (s *APIServer) initRoutes() chi.Router {
	mux := chi.NewMux()

	routes := []Route{
		// User Routes
		{http.MethodPost, "/login", LoginHandler(s.UserService)},
		{http.MethodPost, "/customers", RegisterUserHandler(s.UserService)},
		{http.MethodGet, "/customers/{id}", GetUserHandler(s.UserService)},
		{http.MethodGet, "/cards/{id}", GetCardHandler(s.UserService)},
		{http.MethodGet, "/addresses/{id}", GetAddressHandler(s.UserService)},
		{http.MethodGet, "/customers/{id}/cards", GetUserCardsHandler(s.UserService)},
		{http.MethodGet, "/customers/{id}/addresses", GetUserAddressesHandler(s.UserService)},
		{http.MethodPost, "/customers/{id}/cards", CreateCardHandler(s.UserService)},
		{http.MethodPost, "/customers/{id}/addresses", CreateAddressHandler(s.UserService)},

		// catalogue Routes
		{http.MethodGet, "/catalogue", ListSocksHandler(s.SockLister)},
		{http.MethodGet, "/catalogue/size", CountTagsHandler(s.SockStore)},
		{http.MethodGet, "/catalogue/{id}", GetSocksHandler(s.SockStore)},
		{http.MethodGet, "/tags", TagsHandler(s.SockStore)},
	}

	for _, r := range routes {
		h := r.Handler
		h = WithLogging(s.Logger)(h)
		mux.Method(r.Method, r.Path, ToStdHandler(h))
	}

	mux.Get("/health", ToStdHandler(HealthCheckHandler(s.HealthChecker)).ServeHTTP)

	mux.Handle("/catalogue/images/*", http.StripPrefix(
		"/catalogue/images/", http.FileServer(http.Dir(s.ImagePath)),
	))

	return mux
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
