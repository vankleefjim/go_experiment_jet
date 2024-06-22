package integtest

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	setupDB()
	// if os.Getenv("INTEGRATION_TEST") != "true" {
	// 	return
	// }
	os.Exit(m.Run())
}

func setupDB() {
	ctx := context.Background()

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15.4",
			ExposedPorts: []string{"5432"},
			Env: map[string]string{
				"POSTGRES_USER":     "user1",
				"POSTGRES_PASSWORD": "pg-pw",
				"POSTGRES_DB":       "things",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	failOn(err)

	info, err := dbContainer.Inspect(ctx)
	failOn(err)
	slog.With("info", info).InfoContext(ctx, "some info")
}

func failOn(err error) {
	if err == nil {
		return
	}
	slog.With("err", err.Error()).Error("failed db test")
	os.Exit(1)
}
