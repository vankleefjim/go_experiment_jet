package main

import (
	"github.com/caarlos0/env/v9"
	"jvk.com/things/internal/config"
	"jvk.com/things/internal/server"
)

func main() {
	cfg := config.Server{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	server.New().Run(cfg)
}
