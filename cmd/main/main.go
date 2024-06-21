package main

import (
	"things/internal/config"

	"things/internal/server"

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
