package application

import (
	"context"
	"sync"

	"github.com/nandehu0323/room-monitor/internal/pkg/modules"

	"golang.org/x/sync/errgroup"
)

// Monitor ...
type Monitor struct {
	context context.Context
	mu      *sync.RWMutex
	sensors []modules.Module
}

// NewMonitor ...
func NewMonitor(ctx context.Context) *Monitor {
	return &Monitor{
		context: ctx,
		mu:      new(sync.RWMutex),
	}
}

// Register ...
func (m *Monitor) Register(monitor modules.Module) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sensors = append(m.sensors, monitor)
}

// Run ...
func (m *Monitor) Run() error {
	eg := errgroup.Group{}
	for _, v := range m.sensors {
		v := v
		eg.Go(func() error {
			return v.Watch()
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
