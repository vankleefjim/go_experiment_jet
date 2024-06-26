package integtest

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/vankleefjim/go_experiment_jet/internal/db"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	// TODO run these tests only with this env var?
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
var (
	dbContainer *postgres.PostgresContainer
	dbConn      *sql.DB
	todoDB      *db.TodoDB
)

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

	dbConn, err = dbconn.SQLConnect(dbconn.Config{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Address:  "localhost",
		Port:     mappedPort.Int(),
	})
	failOn(err)

	failOn(migrate.Up(ctx, dbConn))

	todoDB = db.NewTodo(dbConn)
}

func cleanupDB() {
	ctx := context.Background()

	failOn(dbConn.Close())
	failOn(dbContainer.Terminate(ctx))
}

func failOn(err error) {
	if err == nil {
		return
	}
	slog.With("err", err.Error()).Error("failed db test")
	os.Exit(1)
}
