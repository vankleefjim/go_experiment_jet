package migrate

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"

	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func UpCmd(cfg dbconn.Config) *cobra.Command {
	c := &cobra.Command{
		Use: "up",
		Run: func(cmd *cobra.Command, _ []string) {
			err := Up(cmd.Context(), cfg)
			if err != nil {
				slog.With("err", err).ErrorContext(cmd.Context(), "unable to migrate up")
				os.Exit(1)
			}
		},
	}
	return c
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Up(ctx context.Context, cfg dbconn.Config) error {
	conn := must(dbconn.SQLConnect(cfg))
	driver := must(postgres.WithInstance(conn, &postgres.Config{}))

	migrationFiles, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("unable to create iofs: %w", err)
	}
	migrator := must(migrate.NewWithInstance("iofs", migrationFiles, "postgres", driver))

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to migrate: %w", err)
	}
	version, dirty, err := migrator.Version()
	if err != nil {
		return fmt.Errorf("unable to get version: %w", err)
	}
	slog.With("new_version", version, "dirty", dirty).InfoContext(ctx, "successfully migrated")
	return nil
}

func must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
