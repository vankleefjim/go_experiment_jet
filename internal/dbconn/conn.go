package dbconn

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	User     string `env:"DB_USER,required,unset"`
	Password string `env:"DB_PASSWORD,required,unset"`
	Name     string `env:"DB_NAME" envDefault:"things"`
	Address  string `env:"DB_ADDRESS" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
}

func SQLConnect(cfg Config) (*sql.DB, error) {
	return sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.Name))
}
