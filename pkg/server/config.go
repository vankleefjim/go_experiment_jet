package server

type Config struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	HTTPPort int    `env:"HTTP_PORT" envDefault:"8080"`
	GRPCPort int    `env:"GRPC_PORT" envDefault:"8180"`
}
