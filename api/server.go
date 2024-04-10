package api

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

var ErrAPIServerClosed = errors.New("api server closed")

type Server struct {
	Addr   string
	Mux    http.Handler
	Logger *zap.Logger

	httpServer *http.Server
	once       sync.Once
	cancel     func()
}

func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	s.httpServer = &http.Server{Addr: s.Addr, Handler: s.Mux}

	go func() {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			s.Logger.Error("api server stopped with a failure", zap.Error(err))
		}
	}()

	s.Logger.Info("starting api server", zap.String("addr", s.Addr))

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.cancel()
	var err error
	s.once.Do(func() {
		s.Logger.Info("shutting down api server")
		err = s.httpServer.Shutdown(context.Background())
	})
	return err
}
