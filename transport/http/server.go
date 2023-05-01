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

type APIServer struct {
	sockLister    sockLister
	healthChecker healthChecker
	sockStore     domain.SockStore
	imagePath     string
	logger        *zap.Logger

	router *mux.Router
}

func NewAPIServer(
	sockLister sockLister,
	logger *zap.Logger,
	healthChecker healthChecker,
	sockStore domain.SockStore,
	imagePath string,
) *APIServer {
	s := &APIServer{
		sockLister:    sockLister,
		router:        mux.NewRouter(),
		logger:        logger,
		healthChecker: healthChecker,
		sockStore:     sockStore,
		imagePath:     imagePath,
	}

	s.initRoutes()
	return s
}

func (s *APIServer) Start(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

func (s *APIServer) initRoutes() {
	routes := []struct {
		method  string
		path    string
		handler Handler
	}{
		{http.MethodGet, "/health", HealthCheckHandler(s.healthChecker)},
		{http.MethodGet, "/catalogue", ListSocksHandler(s.sockLister)},
		{http.MethodGet, "/catalogue/size", CountTagsHandler(s.sockStore)},
		{http.MethodGet, "/catalogue/{id}", GetSocksHandler(s.sockStore)},
		{http.MethodGet, "/tags", TagsHandler(s.sockStore)},
	}

	for _, route := range routes {
		h := WithLogging(s.logger)(route.handler)
		s.router.Methods(route.method).Path(route.path).Handler(ToStdHandler(h))
	}

	imgServer := http.FileServer(http.Dir(s.imagePath))
	imgServer = http.StripPrefix("/catalogue/images/", imgServer)

	s.router.Methods(http.MethodGet).PathPrefix("/catalogue/images/").Handler(imgServer)
}
