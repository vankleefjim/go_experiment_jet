package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/caarlos0/env/v9"
	"google.golang.org/grpc"
)

type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	done       chan struct{}
}

func New() *Server {
	return &Server{
		done: make(chan struct{}),
	}
}

func (s *Server) Run(
	registerRoutes func(*http.ServeMux),
	registerGRPC func(*grpc.Server), // TODO would be better with option patter but whatever.
) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux)

	httpServer := &http.Server{
		Handler: mux,
		Addr:    httpAddr(cfg),
	}

	s.httpServer = httpServer
	go func() {
		// TODO separate net.Listen out of this like https://go.dev/play/p/mX0OsNWFf-f
		// to make sure it is listening already when starting GRPC server?
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

	if registerGRPC != nil {
		grpcListener, err := net.Listen("tcp", grpcAddr(cfg))
		if err != nil {
			slog.With("err", err).Error("failed listening to net for grpc")
			panic(err)
		}

		grpcServer := grpc.NewServer()
		s.grpcServer = grpcServer
		registerGRPC(grpcServer)
		go func() {
			err := grpcServer.Serve(grpcListener)
			if err != nil {
				if errors.Is(err, grpc.ErrServerStopped) {
					slog.Info("grpc server shut down")
				} else {
					slog.With("err", err).Error("unable to start HTTP server")
					panic(err)
				}
			}
		}()

		slog.With("addr", grpcListener.Addr()).Info("starting GRPC server")
	}

	<-s.done
}

func httpAddr(cfg Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)
}
func grpcAddr(cfg Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort)
}

func (s *Server) Shutdown(ctx context.Context) {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		slog.With("err", err).ErrorContext(ctx, "unable to shutdown HTTP server")
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		// also stops the listener itself
	}

	close(s.done)
}
