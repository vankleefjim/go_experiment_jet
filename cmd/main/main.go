package main

import (
	"github.com/vankleefjim/go_experiment_jet/internal/config"

	"github.com/vankleefjim/go_experiment_jet/internal/server"

	"github.com/caarlos0/env/v9"
)

func main() {
	cfg := config.Server{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	server.New().Run(cfg)
}
