package api

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/oshankkumar/sockshop/api/httpkit"
	"github.com/oshankkumar/sockshop/api/middleware"
	"github.com/oshankkumar/sockshop/api/router"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	Addr          string
	Logger        *zap.Logger
	HealthChecker HealthChecker
	Router        router.Router

	httpServer *http.Server
	once       sync.Once
	cancel     func()
}

func (s *Server) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	mux := chi.NewMux()
	mux.Method(http.MethodGet, "/health", http.HandlerFunc(s.health))
	mux.Method(http.MethodGet, "/metrics", promhttp.Handler())

	middlewareFunc := httpkit.ChainMiddleware(
		middleware.WithLog(s.Logger),
		middleware.WithMetrics(),
		middleware.WithHTTPErrStatus,
	)

	for _, rt := range s.Router.Routes() {
		handler := middlewareFunc(rt.Method, rt.Pattern, rt.Handler)
		mux.Method(rt.Method, rt.Pattern, httpkit.DiscardErr(handler))
	}

	s.httpServer = &http.Server{Addr: s.Addr, Handler: mux}

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

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	hh, err := s.HealthChecker.CheckHealth(r.Context())
	if err != nil {
		httpkit.RespondError(w, &httpkit.Error{Code: http.StatusInternalServerError, Message: err.Error(), Err: err})
		return
	}

	httpkit.RespondJSON(w, HealthResponse{Healths: hh}, http.StatusOK)
}
