package api

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/oshankkumar/sockshop/api/httpkit"
	"github.com/oshankkumar/sockshop/api/middleware"
	"github.com/oshankkumar/sockshop/api/router"

	"go.uber.org/zap"
)

type Server struct {
	Addr          string
	Logger        *zap.Logger
	HealthChecker HealthChecker

	Router     router.Router
	httpServer *http.Server
	once       sync.Once
	cancel     func()
}

func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	s.httpServer = &http.Server{Addr: s.Addr, Handler: s.createMux()}

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

func (s *Server) createMux() router.Mux {
	var mux router.Mux = router.NewMux()
	mux = router.NewInstrumentedMux(mux, middleware.WithLog(s.Logger))

	mux.Method(http.MethodGet, "/health", httpkit.HandlerFunc(s.health))
	s.Router.InstallRoutes(mux)

	return mux
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

func (s *Server) health(w http.ResponseWriter, r *http.Request) *httpkit.Error {
	hh, err := s.HealthChecker.CheckHealth(r.Context())
	if err != nil {
		return &httpkit.Error{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
	}

	httpkit.RespondJSON(w, HealthResponse{Healths: hh}, http.StatusOK)
	return nil
}
