package integtest

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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

func setupDB() {
	ctx := context.Background()
	// dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	// 	ContainerRequest: testcontainers.ContainerRequest{
	// 		Image:        "postgres:15.4",
	// 		ExposedPorts: []string{dbPortStr},
	// 		Env: map[string]string{
	// 			"POSTGRES_USER":     dbUser,
	// 			"PGUSER":            dbUser,
	// 			"POSTGRES_PASSWORD": dbPassword,
	// 			"POSTGRES_DB":       dbName,
	// 		},
	// 		WaitingFor: wait.ForExec([]string{
	// 			"pg_isready", "-U", dbUser, "-d", dbName}),
	// 		// WaitingFor: wait.ForListeningPort(nat.Port(dbPortStr + "/tcp")),
	// 		LogConsumerCfg: &testcontainers.LogConsumerConfig{
	// 			Consumers: []testcontainers.LogConsumer{
	// 				&myOtherLogger{},
	// 			},
	// 		},
	// 		Mounts: testcontainers.ContainerMounts{
	// 			{
	// 				Target: "/var/lib/postgresql/data",
	// 			},
	// 		},
	// 	},
	// 	Started: true,
	// 	Logger:  &myLogger{},
	// })
	// failOn(err)
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.4-alpine"),
		// postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
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
	slog.With("mapped port", mappedPort).InfoContext(ctx, "mapped port?")

	// t.Cleanup(func() {
	// 	if err := pgContainer.Terminate(ctx); err != nil {
	// 		t.Fatalf("failed to terminate pgContainer: %s", err)
	// 	}
	// })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	failOn(err)

	time.Sleep(5 * time.Second) // just to try it?
	slog.With("connstr", connStr).InfoContext(ctx, "connstring found")
	// host, err := dbContainer.Host(ctx)
	// failOn(err)
	// slog.With("host", host).InfoContext(ctx, "host")
	// networks, err := dbContainer.Networks(ctx)
	// failOn(err)
	// slog.With("networks", networks).InfoContext(ctx, "networks")
	// info, err := dbContainer.Inspect(ctx)
	// failOn(err)
	// slog.With("info", info).InfoContext(ctx, "some info")
	// state, err := dbContainer.State(ctx)
	// failOn(err)
	// slog.With("state", state.Status).InfoContext(ctx, "state")
	// _, r, err := dbContainer.Exec(ctx, []string{"pg_isready", "-U", dbUser, "-d", dbName})
	// failOn(err)
	// read, err := io.ReadAll(r)
	// failOn(err)
	// slog.With("result", read).InfoContext(ctx, "exec result")
	failOn(migrate.Up(ctx, dbconn.Config{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Address:  "localhost",
		Port:     mappedPort.Int(),
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
	setupDB()
}

type myLogger struct{}

func (m *myLogger) Printf(format string, v ...interface{}) {
	l := slog.Default()
	for i, v := range v {
		l = l.With(strconv.Itoa(i), v)
	}
	l = l.With("message", format)
	l.Info("printf inside docker")
}

type myOtherLogger struct{}

func (m *myOtherLogger) Accept(l testcontainers.Log) {
	slog.With("content", string(l.Content), "type", l.LogType).Info("msg in container?")
}
