package api

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var ErrAPIServerClosed = errors.New("api server closed")

type Server struct {
	Addr   string
	Mux    http.Handler
	Logger *zap.Logger

	httpServer *http.Server
	cancel     func()
	stopErr    chan error
}

func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	s.httpServer = &http.Server{Addr: s.Addr, Handler: s.Mux}
	s.stopErr = make(chan error, 1)

	go func() {
		<-ctx.Done()
		s.Logger.Info("shutting down api server")
		s.stopErr <- s.httpServer.Shutdown(context.Background())
		close(s.stopErr)
	}()

	s.Logger.Info("starting api server", zap.String("addr", s.Addr))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.cancel()

	if err, ok := <-s.stopErr; ok {
		return err
	}

	return ErrAPIServerClosed
}
