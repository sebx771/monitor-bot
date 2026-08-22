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
		return nil, fmt.Errorf("interval debe ser mayor que cero")
	}

	if task == nil {
		return nil, fmt.Errorf("task no puede ser nil")
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

	// Ejecutar inmediatamente al arrancar la primera vez (Opcional)
	w.execute(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			failed := w.execute(ctx)

			// Si falló y hay un cooldown configurado
			if failed && w.cooldown > 0 {
				timer := time.NewTimer(w.cooldown)

				select {
				case <-ctx.Done():
					timer.Stop() // Evita fuga de memoria
					return nil
				case <-timer.C:
					// Reiniciamos el ticker para que el intervalo de espera
					// vuelva a contar A PARTIR de que terminó el cooldown
					ticker.Reset(w.interval)
				}
			}
		}
	}
}

// execute maneja el bloqueo y la ejecución de la tarea.
// Retorna true si la tarea falló para activar el cooldown.
func (w *Worker) execute(ctx context.Context) bool {
	if !w.mu.TryLock() {
		log.Warn("ciclo anterior aún en ejecución, se omite este ciclo")
		return false
	}
	defer w.mu.Unlock()

	// Si el contexto ya fue cancelado antes de empezar la tarea
	if ctx.Err() != nil {
		return false
	}

	if err := w.task(ctx); err != nil {
		log.Error("ciclo fallido", "error", err)
		return true
	}

	return false
}