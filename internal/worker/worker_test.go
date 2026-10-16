package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunRunsTaskPeriodically(t *testing.T) {
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
		t.Fatalf("expected at least 3 executions, got %d", runs.Load())
	}
}

func TestRunReturnsNilWithCancelledContext(t *testing.T) {
	w, err := New(10*time.Millisecond, 0, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := w.Run(ctx); err != nil {
		t.Fatalf("Run should return nil when the context is cancelled, got: %v", err)
	}
}

func TestRunBackoffAfterFailure(t *testing.T) {
	var runs atomic.Int32

	w, err := New(10*time.Millisecond, 30*time.Millisecond, func(ctx context.Context) error {
		runs.Add(1)
		return errors.New("failure")
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
		t.Fatalf("with a 30ms backoff and a 60ms window, 2 executions were expected, got %d", runs.Load())
	}
}

func TestNewRejectsInvalidParameters(t *testing.T) {
	if _, err := New(0, 0, func(ctx context.Context) error { return nil }); err == nil {
		t.Fatal("expected error with invalid interval")
	}

	if _, err := New(time.Second, 0, nil); err == nil {
		t.Fatal("expected error with nil task")
	}
}

func TestRunManualTrigger(t *testing.T) {
	var runs atomic.Int32

	// Set a very long interval so it only runs on startup and manual triggers
	w, err := New(1*time.Hour, 0, func(ctx context.Context) error {
		runs.Add(1)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	// Wait for the startup execution
	time.Sleep(20 * time.Millisecond)
	if runs.Load() != 1 {
		t.Fatalf("expected 1 execution on startup, got %d", runs.Load())
	}

	// Trigger manual run
	if !w.Trigger() {
		t.Fatal("expected trigger to succeed")
	}

	// Allow execution to complete
	time.Sleep(30 * time.Millisecond)
	if runs.Load() != 2 {
		t.Fatalf("expected 2 executions after trigger, got %d", runs.Load())
	}

	cancel()
	<-done
}

