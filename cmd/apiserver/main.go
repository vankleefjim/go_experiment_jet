package main

import (
	"context"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/api"
	"google.golang.org/grpc"

	"github.com/vankleefjim/go_experiment_jet/pkg/lifecycle"
	"github.com/vankleefjim/go_experiment_jet/pkg/server"

	"github.com/caarlos0/env/v9"
)

func main() {
	cfg := api.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	lcManager := lifecycle.NewManager(context.Background())

	s := server.New()
	lcManager.Register(s)

	apiServer := api.New(cfg)
	lcManager.Register(apiServer)
	s.Run(
		func(mux *http.ServeMux) { apiServer.RegisterRoutes(mux) },
		func(s *grpc.Server) { apiServer.RegisterGRPC(s) },
	)
}
