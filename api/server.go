package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	Addr   string
	Mux    http.Handler
	Logger *zap.Logger

	httpServer *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	s.httpServer = &http.Server{Addr: s.Addr, Handler: s.Mux}
	errc := make(chan error, 1)

	s.Logger.Info("starting app http server :9090")

	go func(errc chan<- error) {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("http server shutdonwn: %w", err)
		}
	}(errc)

	shutdown := func(timeout time.Duration) error {
		s.Logger.Info("received context cancellation; shutting down server")

		shutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := s.httpServer.Shutdown(shutCtx); err != nil {
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
