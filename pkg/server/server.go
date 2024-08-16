package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v9"
)

type Server struct {
	httpServer *http.Server
	done       chan struct{}
}

func New() *Server {
	return &Server{
		done: make(chan struct{}),
	}
}

func (s *Server) Run(registerRoutes func(context.Context, *http.ServeMux)) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	ctx, listenShutdown := setupShutdown(s.Shutdown)
	go listenShutdown()

	mux := http.NewServeMux()
	registerRoutes(ctx, mux)

	httpServer := &http.Server{
		Handler: mux,
		Addr:    addr(cfg),
	}

	s.httpServer = httpServer
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("http server shut down")
			} else {
				slog.With("err", err).Error("unable to start HTTP server")
				panic(err)
			}
		}
	}()
	slog.With("addr", s.httpServer.Addr).Info("starting HTTP server")

	<-s.done
}

func addr(cfg Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func setupShutdown(shutdown func(context.Context)) (notifyCtx context.Context, listenShutdown func()) {
	errC := make(chan error, 1)
	ctx := context.Background()
	notifyCtx, stop := signal.NotifyContext(ctx,
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	return notifyCtx, func() {
		defer func() {
			shutdown(ctx)
			stop()
			close(errC)
		}()
		<-notifyCtx.Done()
		slog.InfoContext(ctx, "signal received, shutting down")
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		slog.With("err", err).ErrorContext(ctx, "unable to shutdown HTTP server")
	}
	close(s.done)
}
