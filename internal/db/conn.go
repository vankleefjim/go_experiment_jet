package db

import (
	"database/sql"
	"fmt"

	"github.com/vankleefjim/go_experiment_jet/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SQLConnect(cfg config.DB) (*sql.DB, error) {
	return sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.Name))
}
