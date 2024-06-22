package integtest

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"
)

func TestMain(m *testing.M) {
	setupDB()
	// if os.Getenv("INTEGRATION_TEST") != "true" {
	// 	return
	// }
	os.Exit(m.Run())
}

const (
	dbUser     = "user1"
	dbPassword = "pg-pw"
	dbName     = "things"
	dbPort     = 5432
)

var (
	dbPortStr = strconv.Itoa(dbPort)
)

func setupDB() {
	ctx := context.Background()

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15.4",
			ExposedPorts: []string{dbPortStr},
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForListeningPort(nat.Port(dbPortStr + "/tcp")),
		},
		Started: true,
	})
	failOn(err)
	time.Sleep(time.Second) // just to try it?

	info, err := dbContainer.Inspect(ctx)
	failOn(err)
	slog.With("info", info).InfoContext(ctx, "some info")

	failOn(migrate.Up(ctx, dbconn.Config{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Address:  "localhost",
		Port:     dbPort,
	}))
}

func failOn(err error) {
	if err == nil {
		return
	}
	slog.With("err", err.Error()).Error("failed db test")
	os.Exit(1)
}

func Test_it(t *testing.T) {
	slog.Info("anything")
}