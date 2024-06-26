package integtest

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	// if os.Getenv("INTEGRATION_TEST") != "true" {
	// 	return
	// }
	setupDB()
	defer cleanupDB()
	os.Exit(m.Run())
}

const (
	dbUser     = "user1"
	dbPassword = "pg-pw"
	dbName     = "things"
)

// Not great but should be fine.
var dbContainer *postgres.PostgresContainer

func setupDB() {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.4-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	failOn(err)
	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	failOn(err)

	failOn(migrate.Up(ctx, dbconn.Config{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Address:  "localhost",
		Port:     mappedPort.Int(),
	}))
}

func cleanupDB() {
	ctx := context.Background()
	err := dbContainer.Terminate(ctx)
	failOn(err)
}

func failOn(err error) {
	if err == nil {
		return
	}
	slog.With("err", err.Error()).Error("failed db test")
	os.Exit(1)
}

func Test_it(t *testing.T) {
}
