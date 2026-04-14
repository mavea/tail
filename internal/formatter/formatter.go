package formatter

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	formatterBuffer "tail/internal/formatter/buffer"
	formatterWindow "tail/internal/formatter/window"
)

type formatter struct {
	cfg Cfg

	window    *formatterWindow.Window
	buffer    buffers
	stderr    buffers
	stderrSet bool
	hardErr   bool
	out       Out
	print     int8
	wg        *sync.WaitGroup
	status    int
	mu        sync.Mutex
}

type ExitError struct {
	Code int
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

func NewFormatter(ctx context.Context, out Out, cfg Cfg, window *formatterWindow.Window) (Render, Cancel, error) {
	result := &formatter{
		cfg:    cfg,
		buffer: formatterBuffer.New(cfg, window),
		stderr: formatterBuffer.New(cfg, window),
		out:    out,
		wg:     &sync.WaitGroup{},
		window: window,

		status: 0,
	}
	ctxLocal, cancel := context.WithCancel(ctx)

	result.buffer.SetDefaultStyle(result.out.GetDefaultStyle())

	if cfg.GetOutputMode() == "thread" {
		result.print = 1

		result.wg.Add(1)
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer func() {
				ticker.Stop()
				result.wg.Done()
			}()
			for {
				select {
				case <-ctxLocal.Done():
					return
				case <-ticker.C:
					_ = result.out.Print()
				}
			}
		}()
	}

	return result, func() error {
		cancel()
		return result.cancel()
	}, nil
}

func (r *formatter) Set(data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buffer.Add(data)
	r.out.SetData(r.buffer.GetLast(r.cfg.GetMaxLineCount()))

	if r.print == 0 {
		return r.out.Print()
	}

	return nil
}

func (r *formatter) SetErrLine(data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stderrSet = true
	r.stderr.Add(data)

	return nil
}

func (r *formatter) SetErr(err error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stderrSet = true
	r.hardErr = true
	r.stderr.Add(err.Error())
	if r.status == 0 {
		r.status = 1
	}

	return nil
}

func (r *formatter) SetStatus(status int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = status

	return nil
}

func (r *formatter) cancel() error {
	r.wg.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	var cancelErr error
	addErr := func(err error) {
		if err != nil {
			cancelErr = errors.Join(cancelErr, err)
		}
	}

	addErr(r.out.ClearScreen())
	if r.stderrSet && (r.hardErr || r.status != 0) {
		errs := r.stderr.GetFull()
		if len(errs) > 0 {
			addErr(r.out.Error(r.buffer.GetFull(), errs))
		}
	} else if r.cfg.IsFullOutput() {
		addErr(r.out.PrintFull(r.buffer.GetFull()))
	}

	if r.status != 0 {
		addErr(r.out.SetStatus(r.status))
		exitErr := ExitError{Code: r.status}
		if cancelErr != nil {
			return errors.Join(exitErr, cancelErr)
		}
		return exitErr
	}

	return cancelErr
}
