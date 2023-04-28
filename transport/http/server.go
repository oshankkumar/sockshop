package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oshankkumar/sockshop/domain"

	"go.uber.org/zap"
)

type Server struct {
	sockLister    sockLister
	router        *mux.Router
	logger        *zap.Logger
	healthChecker healthChecker
	sockStore     domain.SockStore
}

func NewServer(
	sockLister sockLister,
	logger *zap.Logger,
	healthChecker healthChecker,
	sockStore domain.SockStore,
) *Server {
	s := &Server{
		sockLister:    sockLister,
		router:        mux.NewRouter(),
		logger:        logger,
		healthChecker: healthChecker,
		sockStore:     sockStore,
	}
	return s
}

func (s *Server) Start(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.initRoutes()

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	errc := make(chan error, 1)

	go func(errc chan<- error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		s.logger.Info("received context cancellation; shutting down server")

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

func (s *Server) initRoutes() {
	routes := []struct {
		method  string
		path    string
		handler HTTPHandler
	}{
		{http.MethodGet, "/health", HealthCheckHandler(s.healthChecker)},
		{http.MethodGet, "/catalogue", ListSocksHandler(s.sockLister)},
		{http.MethodGet, "/catalogue/size", CountTagsHandler(s.sockStore)},
		{http.MethodGet, "/catalogue/{id}", GetSockHandler(s.sockStore)},
	}

	for _, route := range routes {
		h := WithLogging(s.logger)(route.handler)
		s.router.Methods(route.method).Path(route.path).Handler(ToStdHandler(h))
	}
}
