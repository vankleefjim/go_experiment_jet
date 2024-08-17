package lifecycle

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// NewManager handles shutdown lifecycles.
func NewManager(ctx context.Context) *manager {
	errC := make(chan error, 1) // TODO maybe remove
	ctx, stop := signal.NotifyContext(ctx,
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	m := &manager{}

	go func() {
		defer func() {
			m.shutdown(ctx)
			stop()
			close(errC)
		}()
		<-ctx.Done()
		slog.InfoContext(ctx, "signal received, shutting down")
	}()

	return m
}

type manager struct {
	mu          sync.Mutex
	shutdowners []Shutdowner
}

type Shutdowner interface {
	Shutdown(context.Context)
}

func (m *manager) Register(s Shutdowner) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shutdowners = append(m.shutdowners, s)
}

func (m *manager) shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.shutdowners {
		s.Shutdown(ctx)
	}
}
