package main

import (
	"context"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/vankleefjim/go_experiment_jet/internal/ui"
	"github.com/vankleefjim/go_experiment_jet/pkg/lifecycle"
	"github.com/vankleefjim/go_experiment_jet/pkg/server"
)

func main() {
	cfg := ui.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	lcManager := lifecycle.NewManager(context.Background())

	s := server.New()
	lcManager.Register(s)

	uiServer := ui.New(cfg)
	s.Run(
		func(mux *http.ServeMux) { uiServer.RegisterRoutes(mux) },
		nil,
	)
}
