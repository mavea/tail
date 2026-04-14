package context

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewMainContext(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := NewMainContext(context.Background(), os.Stdout, &wg)
	if ctx == nil || cancel == nil {
		t.Fatal("expected context and cancel")
	}
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("context cleanup timed out")
	}
}

func TestNewMainContextCancelation(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := NewMainContext(context.Background(), os.Stdout, &wg)
	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("context should be done after cancel")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("waitgroup should be released")
	}
}
