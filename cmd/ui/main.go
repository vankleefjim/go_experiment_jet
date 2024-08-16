package main

import (
	"context"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/vankleefjim/go_experiment_jet/internal/ui"
	"github.com/vankleefjim/go_experiment_jet/pkg/server"
)

func main() {
	cfg := ui.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	server.New().Run(func(ctx context.Context, mux *http.ServeMux) { ui.RegisterRoutes(ctx, cfg, mux) })
}
