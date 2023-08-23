package config

type Server struct {
	DB   DB
	HTTP HTTP
}

type HTTP struct {
	Host string `env:"HTTP_HOST" envDefault:"localhost"`
	Port int    `env:"HTTP_PORT" envDefault:"8080"`
}

type DB struct {
	User     string `env:"DB_USER,required,unset"`
	Password string `env:"DB_PASSWORD,required,unset"`
	Name     string `env:"DB_NAME" envDefault:"things"`
	Address  string `env:"DB_ADDRESS" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
}
