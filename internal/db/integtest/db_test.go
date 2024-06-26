package integtest

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"
)

func TestMain(m *testing.M) {
	//setupDB()
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

func setupDB(t *testing.T) {
	ctx := context.Background()

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15.4",
			ExposedPorts: []string{dbPortStr},
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"PGUSER":            dbUser,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForExec([]string{
				"pg_isready", "-U", dbUser, "-d", dbName}),
			// WaitingFor: wait.ForListeningPort(nat.Port(dbPortStr + "/tcp")),
		},
		Started: true,
		Logger:  testcontainers.TestLogger(t),
	})
	failOn(err)
	time.Sleep(10 * time.Second) // just to try it?
	host, err := dbContainer.Host(ctx)
	failOn(err)
	slog.With("host", host).InfoContext(ctx, "host")
	networks, err := dbContainer.Networks(ctx)
	failOn(err)
	slog.With("networks", networks).InfoContext(ctx, "networks")
	info, err := dbContainer.Inspect(ctx)
	failOn(err)
	slog.With("info", info).InfoContext(ctx, "some info")
	_, r, err := dbContainer.Exec(ctx, []string{"pg_isready", "-U", dbUser, "-d", dbName})
	failOn(err)
	read, err := io.ReadAll(r)
	failOn(err)
	t.Error(string(read))
	execR, err := exec.Command("nc -zv localhost 5432").Output()
	failOn(err)
	t.Error(string(execR))
	failOn(migrate.Up(ctx, dbconn.Config{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Address:  host,
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
	setupDB(t)
}
