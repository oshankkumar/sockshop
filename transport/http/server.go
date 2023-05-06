package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/domain"

	"go.uber.org/zap"
)

type APIServer struct {
	HealthChecker HealthChecker
	UserService   api.UserService
	SockLister    SockLister
	SockStore     domain.SockStore
	ImagePath     string
	Logger        *zap.Logger
}

func (s *APIServer) Start(ctx context.Context, addr string) error {
	router := mux.NewRouter()
	s.installRoutes(router)

	srv := &http.Server{Addr: addr, Handler: router}

	errc := make(chan error, 1)

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

func (s *APIServer) installRoutes(router *mux.Router) {
	routes := []struct {
		method  string
		path    string
		handler Handler
	}{
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

		// Catalogue routes
		{http.MethodGet, "/health", HealthCheckHandler(s.HealthChecker)},
		{http.MethodGet, "/catalogue", ListSocksHandler(s.SockLister)},
		{http.MethodGet, "/catalogue/size", CountTagsHandler(s.SockStore)},
		{http.MethodGet, "/catalogue/{id}", GetSocksHandler(s.SockStore)},
		{http.MethodGet, "/tags", TagsHandler(s.SockStore)},
	}

	for _, route := range routes {
		h := WithLogging(s.Logger)(route.handler)
		router.Methods(route.method).Path(route.path).Handler(ToStdHandler(h))
	}

	imgServer := http.FileServer(http.Dir(s.ImagePath))
	imgServer = http.StripPrefix("/catalogue/images/", imgServer)

	router.Methods(http.MethodGet).PathPrefix("/catalogue/images/").Handler(imgServer)
}
