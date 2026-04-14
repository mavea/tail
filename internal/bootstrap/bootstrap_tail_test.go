package bootstrap

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"tail/internal/formatter"
	formatterWindow "tail/internal/formatter/window"
	sinkDomain "tail/internal/sink/domain"
	sourceGeneral "tail/internal/source/general"
	"tail/tests/mocks"
)

func TestTailSuccessWithInjectedFactories(t *testing.T) {
	origScanner := newScanner
	origWriter := newWriter
	origFormatter := newFormatter
	defer func() {
		newScanner = origScanner
		newWriter = origWriter
		newFormatter = origFormatter
	}()

	outCh := make(chan string, 1)
	outCh <- "line"
	close(outCh)
	errCh := make(chan string)
	close(errCh)

	scanner := &mocks.ScannerMock{
		OutFunc:       func() <-chan string { return outCh },
		ErrFunc:       func() <-chan string { return errCh },
		GetStatusFunc: func() (int, error) { return 0, nil },
	}

	render := &mocks.RenderMock{
		SetFunc:        func(string) error { return nil },
		SetErrLineFunc: func(string) error { return nil },
		SetErrFunc:     func(error) error { return nil },
		SetStatusFunc:  func(int) error { return nil },
	}

	outTarget := &mocks.OutMock{
		GetDefaultStyleFunc: func() (string, string, string) { return "", "", "" },
		PrintFunc:           func() error { return nil },
		SetDataFunc:         func([]string) {},
		ClearScreenFunc:     func() error { return nil },
		ErrorFunc:           func([]string, []string) error { return nil },
		SetStatusFunc:       func(int) error { return nil },
	}

	newScanner = func(ctx context.Context, cfg Cfg) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
		return scanner, func() error { return nil }, nil
	}
	newWriter = func(context.Context, Cfg, *bufio.Writer, *formatterWindow.Window) (sinkDomain.Target, error) {
		return outTarget, nil
	}
	newFormatter = func(ctx context.Context, out formatter.Out, cfg formatter.Cfg, window *formatterWindow.Window) (formatter.Render, formatter.Cancel, error) {
		return render, func() error { return nil }, nil
	}

	var b bytes.Buffer
	rm := NewRunMode(testBootstrapCfg{}, nil, os.Stdin, bufio.NewWriter(&b), bufio.NewWriter(&b))
	if err := rm.tail(context.Background()); err != nil {
		t.Fatalf("tail returned error: %v", err)
	}
}

func TestTailScannerError(t *testing.T) {
	origScanner := newScanner
	defer func() { newScanner = origScanner }()

	newScanner = func(context.Context, Cfg) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
		return nil, nil, errors.New("scan fail")
	}

	var b bytes.Buffer
	rm := NewRunMode(testBootstrapCfg{}, nil, os.Stdin, bufio.NewWriter(&b), bufio.NewWriter(&b))
	if err := rm.tail(context.Background()); err == nil {
		t.Fatal("expected tail to fail")
	}
}

func TestTailWriterError(t *testing.T) {
	origScanner := newScanner
	origWriter := newWriter
	defer func() {
		newScanner = origScanner
		newWriter = origWriter
	}()

	newScanner = func(context.Context, Cfg) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
		s := &mocks.ScannerMock{
			OutFunc:       func() <-chan string { ch := make(chan string); close(ch); return ch },
			ErrFunc:       func() <-chan string { ch := make(chan string); close(ch); return ch },
			GetStatusFunc: func() (int, error) { return 0, nil },
		}
		return s, func() error { return nil }, nil
	}
	newWriter = func(context.Context, Cfg, *bufio.Writer, *formatterWindow.Window) (sinkDomain.Target, error) {
		return nil, errors.New("writer failed")
	}

	var b bytes.Buffer
	rm := NewRunMode(testBootstrapCfg{}, nil, os.Stdin, bufio.NewWriter(&b), bufio.NewWriter(&b))
	if err := rm.tail(context.Background()); err == nil {
		t.Fatal("expected writer error")
	}
}

func TestTailFormatterError(t *testing.T) {
	origScanner := newScanner
	origWriter := newWriter
	origFormatter := newFormatter
	defer func() {
		newScanner = origScanner
		newWriter = origWriter
		newFormatter = origFormatter
	}()

	newScanner = func(context.Context, Cfg) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
		s := &mocks.ScannerMock{
			OutFunc:       func() <-chan string { ch := make(chan string); close(ch); return ch },
			ErrFunc:       func() <-chan string { ch := make(chan string); close(ch); return ch },
			GetStatusFunc: func() (int, error) { return 0, nil },
		}
		return s, func() error { return nil }, nil
	}
	newWriter = func(context.Context, Cfg, *bufio.Writer, *formatterWindow.Window) (sinkDomain.Target, error) {
		return &mocks.OutMock{
			GetDefaultStyleFunc: func() (string, string, string) { return "", "", "" },
			PrintFunc:           func() error { return nil },
			SetDataFunc:         func([]string) {},
			ClearScreenFunc:     func() error { return nil },
			ErrorFunc:           func([]string, []string) error { return nil },
			SetStatusFunc:       func(int) error { return nil },
		}, nil
	}
	newFormatter = func(context.Context, formatter.Out, formatter.Cfg, *formatterWindow.Window) (formatter.Render, formatter.Cancel, error) {
		return nil, nil, errors.New("formatter failed")
	}

	var b bytes.Buffer
	rm := NewRunMode(testBootstrapCfg{}, nil, os.Stdin, bufio.NewWriter(&b), bufio.NewWriter(&b))
	if err := rm.tail(context.Background()); err == nil {
		t.Fatal("expected formatter creation error")
	}
}

func TestTailPipeAndCancelErrors(t *testing.T) {
	origScanner := newScanner
	origWriter := newWriter
	origFormatter := newFormatter
	defer func() {
		newScanner = origScanner
		newWriter = origWriter
		newFormatter = origFormatter
	}()

	outCh := make(chan string, 1)
	outCh <- "line"
	close(outCh)
	errCh := make(chan string)
	close(errCh)

	newScanner = func(context.Context, Cfg) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
		s := &mocks.ScannerMock{
			OutFunc:       func() <-chan string { return outCh },
			ErrFunc:       func() <-chan string { return errCh },
			GetStatusFunc: func() (int, error) { return 0, nil },
		}
		return s, func() error { return errors.New("scanner cancel failed") }, nil
	}
	newWriter = func(context.Context, Cfg, *bufio.Writer, *formatterWindow.Window) (sinkDomain.Target, error) {
		return &mocks.OutMock{
			GetDefaultStyleFunc: func() (string, string, string) { return "", "", "" },
			PrintFunc:           func() error { return nil },
			SetDataFunc:         func([]string) {},
			ClearScreenFunc:     func() error { return nil },
			ErrorFunc:           func([]string, []string) error { return nil },
			SetStatusFunc:       func(int) error { return nil },
		}, nil
	}
	newFormatter = func(context.Context, formatter.Out, formatter.Cfg, *formatterWindow.Window) (formatter.Render, formatter.Cancel, error) {
		return &mocks.RenderMock{
			SetFunc:        func(string) error { return errors.New("set failed") },
			SetErrLineFunc: func(string) error { return nil },
			SetErrFunc:     func(error) error { return nil },
			SetStatusFunc:  func(int) error { return nil },
		}, func() error { return errors.New("formatter cancel failed") }, nil
	}

	var b bytes.Buffer
	rm := NewRunMode(testBootstrapCfg{}, nil, os.Stdin, bufio.NewWriter(&b), bufio.NewWriter(&b))
	if err := rm.tail(context.Background()); err == nil {
		t.Fatal("expected tail error from pipe/cancel branches")
	}
}
