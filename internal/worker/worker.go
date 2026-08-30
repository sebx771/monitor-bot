package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sebx771/monitor-bot/internal/logger"
)

var log = logger.NewLogger("WORKER")

type Task func(ctx context.Context) error

type Worker struct {
	interval time.Duration
	cooldown time.Duration
	task     Task

	mu sync.Mutex
}

func New(interval, cooldown time.Duration, task Task) (*Worker, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than zero")
	}

	if task == nil {
		return nil, fmt.Errorf("task cannot be nil")
	}

	return &Worker{
		interval: interval,
		cooldown: cooldown,
		task:     task,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run immediately on startup the first time (optional)
	w.execute(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			failed := w.execute(ctx)

			// If it failed and a cooldown is configured
			if failed && w.cooldown > 0 {
				timer := time.NewTimer(w.cooldown)

				select {
				case <-ctx.Done():
					timer.Stop() // Avoids memory leak
					return nil
				case <-timer.C:
					// Reset the ticker so the wait interval counts
					// again FROM when the cooldown finished
					ticker.Reset(w.interval)
				}
			}
		}
	}
}

// execute handles the locking and execution of the task.
// Returns true if the task failed to trigger the cooldown.
func (w *Worker) execute(ctx context.Context) bool {
	if !w.mu.TryLock() {
		log.Warn("previous cycle still running, skipping this cycle")
		return false
	}
	defer w.mu.Unlock()

	// If the context was already cancelled before starting the task
	if ctx.Err() != nil {
		return false
	}

	if err := w.task(ctx); err != nil {
		log.Error("cycle failed", "error", err)
		return true
	}

	return false
}