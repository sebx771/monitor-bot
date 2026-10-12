package aiven

import (
	"context"
	"errors"
	"fmt"

	config "github.com/sebx771/monitor-bot/internal/config"
	"github.com/sebx771/monitor-bot/internal/status"
	"github.com/sebx771/monitor-bot/internal/worker"
)

// NewTask builds the worker task for checking Aiven services across all credentials.
func NewTask(credentials []config.Credential, store *status.Store) worker.Task {
	return func(ctx context.Context) error {
		var allServices []status.MonitoredService
		var errs []error
		successfulProjects := 0

		for _, cred := range credentials {
			checker := NewChecker(NewClient(cred.Token), cred.Project)

			services, err := checker.Check()
			if err != nil {
				errs = append(errs, fmt.Errorf("project %s: %w", cred.Project, err))
				continue
			}
			successfulProjects++

			for _, s := range services {
				allServices = append(allServices, status.MonitoredService{
					Target: s.Name,
					Scope:  cred.Project,
					State:  s.State,
				})
			}
		}

		var ws status.WorkerStatus

		switch {
		case len(errs) == 0:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Operational, allServices, "")
		case successfulProjects == 0:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Error, allServices, "All projects failed")
		default:
			ws, _ = status.NewWorkerStatus("aiven", "Aiven", status.Degraded, allServices, "Some projects failed to check")
		}

		if store != nil {
			store.Set(ws)
		}

		if len(errs) > 0 && successfulProjects == 0 {
			return errors.Join(errs...)
		}

		return nil
	}
}
