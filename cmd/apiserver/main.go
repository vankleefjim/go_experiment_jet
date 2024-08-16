package main

import (
	"context"
	"net/http"

	"github.com/vankleefjim/go_experiment_jet/internal/api"

	"github.com/vankleefjim/go_experiment_jet/pkg/server"

	"github.com/caarlos0/env/v9"
)

func main() {
	cfg := api.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	server.New().Run(func(ctx context.Context, mux *http.ServeMux) { api.RegisterRoutes(ctx, cfg, mux) })
}
