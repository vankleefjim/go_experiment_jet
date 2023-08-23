package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"jvk.com/things/internal/config"
)

func SQLConnect(cfg config.DB) (*sql.DB, error) {
	return sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.Name))
}
