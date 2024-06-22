package config

import "github.com/vankleefjim/go_experiment_jet/pkg/dbconn"

type Server struct {
	DB   dbconn.Config
	HTTP HTTP
}

type HTTP struct {
	Host string `env:"HTTP_HOST" envDefault:"localhost"`
	Port int    `env:"HTTP_PORT" envDefault:"8080"`
}
