package migrate

import (
	"errors"
	"log/slog"
	"os"

	"github.com/vankleefjim/go_experiment_jet/internal/config"
	"github.com/vankleefjim/go_experiment_jet/internal/db"

	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(cfg config.DB) *cobra.Command {
	c := &cobra.Command{
		Use: "up",
		Run: func(cmd *cobra.Command, args []string) {
			conn := must(db.SQLConnect(cfg))
			driver := must(postgres.WithInstance(conn, &postgres.Config{}))
			migrator := must(migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver))

			if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				slog.With("err", err).ErrorContext(cmd.Context(), "unable to migrate")
				os.Exit(1)
			}
			version, dirty, err := migrator.Version()
			if err != nil {
				slog.With("err", err).ErrorContext(cmd.Context(), "unable to get version")
				os.Exit(1)
			}
			slog.With("new_version", version, "dirty", dirty).InfoContext(cmd.Context(), "successfully migrated")
		},
	}
	return c
}
