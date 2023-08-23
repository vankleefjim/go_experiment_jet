package api

import (
	"database/sql"
	"net/http"

	"golang.org/x/exp/slog"
	"jvk.com/things/internal/config"
	"jvk.com/things/internal/db"
	"jvk.com/things/internal/todos"
)

func Routes(cfg config.Server, conn *sql.DB) *http.ServeMux {
	// This extra level is not needed when writing a microservice
	mux := http.NewServeMux()

	// TODO things like CORS

	mux.HandleFunc("/ping", pong)

	todoDB := db.NewTodos(conn)
	mux.HandleFunc("/todo", todos.New(todoDB).Routes())
	return mux
}

func pong(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("PONG")); err != nil {
		slog.With("err", err).Error("unable to respond to ping")
	}
}
