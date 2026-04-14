package bootstrap

import (
	"context"
	"errors"
	"testing"

	mocks2 "tail/tests/mocks"
)

func TestPipeUsesGeneratedMocksHappyPath(t *testing.T) {
	outCh := make(chan string, 2)
	outCh <- "one"
	outCh <- "two"
	close(outCh)

	errCh := make(chan string)
	close(errCh)

	scanner := &mocks2.ScannerMock{
		OutFunc:       func() <-chan string { return outCh },
		ErrFunc:       func() <-chan string { return errCh },
		GetStatusFunc: func() (int, error) { return 0, nil },
	}

	render := &mocks2.RenderMock{
		SetFunc:        func(string) error { return nil },
		SetErrLineFunc: func(string) error { return nil },
		SetErrFunc:     func(error) error { return nil },
		SetStatusFunc:  func(int) error { return nil },
	}

	rm := &RunMode{}
	if err := rm.pipe(context.Background(), scanner, render); err != nil {
		t.Fatalf("unexpected pipe error: %v", err)
	}

	if len(render.SetCalls()) != 2 {
		t.Fatalf("expected 2 Set calls, got %d", len(render.SetCalls()))
	}
	if len(render.SetErrCalls()) != 0 {
		t.Fatalf("expected no SetErr calls, got %d", len(render.SetErrCalls()))
	}
	if len(render.SetStatusCalls()) != 1 || render.SetStatusCalls()[0].Status != 0 {
		t.Fatalf("expected SetStatus(0), got %+v", render.SetStatusCalls())
	}
}

func TestPipeUsesGeneratedMocksErrChannel(t *testing.T) {
	outCh := make(chan string)
	close(outCh)

	errCh := make(chan string, 1)
	errCh <- "boom"
	close(errCh)

	scanner := &mocks2.ScannerMock{
		OutFunc:       func() <-chan string { return outCh },
		ErrFunc:       func() <-chan string { return errCh },
		GetStatusFunc: func() (int, error) { return 0, nil },
	}

	render := &mocks2.RenderMock{
		SetFunc:        func(string) error { return nil },
		SetErrLineFunc: func(string) error { return nil },
		SetErrFunc:     func(error) error { return nil },
		SetStatusFunc:  func(int) error { return nil },
	}

	rm := &RunMode{}
	if err := rm.pipe(context.Background(), scanner, render); err != nil {
		t.Fatalf("unexpected pipe error: %v", err)
	}

	if len(render.SetErrLineCalls()) != 1 {
		t.Fatalf("expected 1 SetErrLine call, got %d", len(render.SetErrLineCalls()))
	}
	if len(render.SetStatusCalls()) != 1 || render.SetStatusCalls()[0].Status != 0 {
		t.Fatalf("expected SetStatus(0), got %+v", render.SetStatusCalls())
	}
}

func TestPipeCancelSetsStatus130(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	scanner := &mocks2.ScannerMock{
		OutFunc:       func() <-chan string { return make(chan string) },
		ErrFunc:       func() <-chan string { return make(chan string) },
		GetStatusFunc: func() (int, error) { return 0, nil },
	}

	render := &mocks2.RenderMock{
		SetFunc:        func(string) error { return nil },
		SetErrLineFunc: func(string) error { return nil },
		SetErrFunc:     func(error) error { return nil },
		SetStatusFunc: func(status int) error {
			if status != 130 {
				t.Fatalf("expected status 130, got %d", status)
			}
			return nil
		},
	}

	rm := &RunMode{}
	if err := rm.pipe(ctx, scanner, render); err != nil {
		t.Fatalf("unexpected pipe error: %v", err)
	}

	if len(render.SetStatusCalls()) != 1 {
		t.Fatalf("expected SetStatus to be called once")
	}
}

func TestPipeGetStatusErrorGoesToSetErr(t *testing.T) {
	statusErr := errors.New("status failed")

	outCh := make(chan string)
	close(outCh)
	errCh := make(chan string)
	close(errCh)

	scanner := &mocks2.ScannerMock{
		OutFunc:       func() <-chan string { return outCh },
		ErrFunc:       func() <-chan string { return errCh },
		GetStatusFunc: func() (int, error) { return 0, statusErr },
	}

	render := &mocks2.RenderMock{
		SetFunc:        func(string) error { return nil },
		SetErrLineFunc: func(string) error { return nil },
		SetStatusFunc:  func(int) error { return nil },
		SetErrFunc: func(err error) error {
			if !errors.Is(err, statusErr) {
				t.Fatalf("unexpected SetErr argument: %v", err)
			}
			return nil
		},
	}

	rm := &RunMode{}
	if err := rm.pipe(context.Background(), scanner, render); err != nil {
		t.Fatalf("unexpected pipe error: %v", err)
	}

	if len(render.SetErrCalls()) != 1 {
		t.Fatalf("expected one SetErr call, got %d", len(render.SetErrCalls()))
	}
}
