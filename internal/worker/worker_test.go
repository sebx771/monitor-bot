package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunEjecutaTareaPeriodicamente(t *testing.T) {
	var runs atomic.Int32

	w, err := New(10*time.Millisecond, 0, func(ctx context.Context) error {
		runs.Add(1)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	time.Sleep(80 * time.Millisecond)
	cancel()

	if err := <-done; err != nil {
		t.Fatal(err)
	}

	if runs.Load() < 3 {
		t.Fatalf("se esperaban al menos 3 ejecuciones, se obtuvieron %d", runs.Load())
	}
}

func TestRunRetornaNilConContextoCancelado(t *testing.T) {
	w, err := New(10*time.Millisecond, 0, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := w.Run(ctx); err != nil {
		t.Fatalf("Run debería retornar nil al cancelar el contexto, obtuvo: %v", err)
	}
}

func TestRunBackoffTrasFallo(t *testing.T) {
	var runs atomic.Int32

	w, err := New(10*time.Millisecond, 30*time.Millisecond, func(ctx context.Context) error {
		runs.Add(1)
		return errors.New("fallo")
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	time.Sleep(60 * time.Millisecond)
	cancel()

	if err := <-done; err != nil {
		t.Fatal(err)
	}

	if runs.Load() != 2 {
		t.Fatalf("con backoff de 30ms y ventana de 60ms se esperaban 2 ejecuciones, se obtuvieron %d", runs.Load())
	}
}

func TestNewRechazaParametrosInvalidos(t *testing.T) {
	if _, err := New(0, 0, func(ctx context.Context) error { return nil }); err == nil {
		t.Fatal("se esperaba error con interval inválido")
	}

	if _, err := New(time.Second, 0, nil); err == nil {
		t.Fatal("se esperaba error con task nil")
	}
}
