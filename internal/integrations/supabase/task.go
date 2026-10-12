package supabase

import (
	"context"

	"github.com/sebx771/monitor-bot/internal/logger"
	"github.com/sebx771/monitor-bot/internal/status"
	"github.com/sebx771/monitor-bot/internal/worker"
)

// NewTask builds the worker task for checking Supabase projects and reporting status.
func NewTask(token string, store *status.Store) worker.Task {
	return func(ctx context.Context) error {
		log := logger.NewLogger("SUPABASE")
		checker := NewChecker(NewClient(token), log)

		projects, err := checker.Check()
		var services []status.MonitoredService
		var wsStatus status.Status
		var msg string

		if err != nil {
			wsStatus = status.Error
			msg = err.Error()
		} else {
			for _, p := range projects {
				services = append(services, status.MonitoredService{
					Target: p.Ref,
					State:  p.Status,
				})
			}
			wsStatus = status.Operational
		}

		ws, _ := status.NewWorkerStatus("supabase", "Supabase", wsStatus, services, msg)
		if store != nil {
			store.Set(ws)
		}

		return err
	}
}
