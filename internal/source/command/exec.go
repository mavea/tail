package command

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	sourceGeneral "tail/internal/source/general"
)

const chanCache = 1
const maxScanTokenSize = 1024 * 1024

type exe struct {
	ctx    context.Context
	cancel context.CancelFunc

	cmd        *exec.Cmd
	outScanner *bufio.Scanner
	outChan    chan string
	errScanner *bufio.Scanner
	errChan    chan string

	outDone chan struct{}
	errDone chan struct{}

	statusOnce sync.Once
	exitCode   int
	statusErr  error
	scanErr    error
	scanErrMu  sync.Mutex
}

func shouldUseShellFallback(programPath string) bool {
	if strings.ContainsAny(programPath, " \t\n") {
		return true
	}

	// Shell meta characters: if present, user likely expects shell expression execution.
	return strings.ContainsAny(programPath, `|&;<>()$\"'*!?[]{}=`)
}

func resolveCommand(programPath string, args []string) (string, []string, error) {
	if strings.TrimSpace(programPath) == "" {
		return "", nil, fmt.Errorf("command is empty")
	}

	if len(args) > 0 {
		return programPath, args, nil
	}

	if _, err := exec.LookPath(programPath); err == nil {
		return programPath, args, nil
	}

	if !shouldUseShellFallback(programPath) {
		return "", nil, fmt.Errorf("failed to resolve executable: %s", programPath)
	}

	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", programPath}, nil
	}

	return "sh", []string{"-c", programPath}, nil
}

func New(ctx context.Context, programPath string, args ...string) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
	binary, resolvedArgs, err := resolveCommand(programPath, args)
	if err != nil {
		return nil, nil, err
	}

	s := &exe{
		outDone: make(chan struct{}),
		errDone: make(chan struct{}),
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	// #nosec G204 -- Command mode intentionally executes user-provided program and args from CLI flags.
	s.cmd = exec.CommandContext(s.ctx, binary, resolvedArgs...)

	stdoutPipe, err := s.cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	s.outScanner = bufio.NewScanner(stdoutPipe)
	s.outScanner.Buffer(make([]byte, 64*1024), maxScanTokenSize)
	s.outChan = make(chan string, chanCache)

	stderrPipe, err := s.cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	s.errScanner = bufio.NewScanner(stderrPipe)
	s.errScanner.Buffer(make([]byte, 64*1024), maxScanTokenSize)
	s.errChan = make(chan string, chanCache)

	if err = s.cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start command: %w", err)
	}
	s.scan()

	return s, s.cancelFunc, nil
}

func (s *exe) scan() {
	go func() {
		defer func() {
			close(s.outChan)
			close(s.outDone)
		}()
		for s.outScanner.Scan() {
			select {
			case s.outChan <- s.outScanner.Text():
			case <-s.ctx.Done():
				return
			}
		}
		if err := s.outScanner.Err(); err != nil {
			s.setScanErr(fmt.Errorf("stdout scanner error: %w", err))
		}
	}()

	go func() {
		defer func() {
			close(s.errChan)
			close(s.errDone)
		}()
		for s.errScanner.Scan() {
			select {
			case s.errChan <- s.errScanner.Text():
			case <-s.ctx.Done():
				return
			}
		}
		if err := s.errScanner.Err(); err != nil {
			s.setScanErr(fmt.Errorf("stderr scanner error: %w", err))
		}
	}()
}

func (s *exe) setScanErr(err error) {
	if err == nil {
		return
	}
	s.scanErrMu.Lock()
	defer s.scanErrMu.Unlock()
	s.scanErr = errors.Join(s.scanErr, err)
}

func (s *exe) Out() <-chan string {
	return s.outChan
}

func (s *exe) Err() <-chan string {
	return s.errChan
}

func (s *exe) GetStatus() (int, error) {
	s.statusOnce.Do(func() {
		err := s.cmd.Wait()
		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				s.exitCode = exitErr.ExitCode()
				s.statusErr = nil
			} else {
				s.exitCode = 1
				s.statusErr = err
			}
		} else {
			s.exitCode = 0
			s.statusErr = nil
		}

		s.scanErrMu.Lock()
		scanErr := s.scanErr
		s.scanErrMu.Unlock()
		if scanErr != nil {
			if s.statusErr == nil {
				s.statusErr = scanErr
			}
			if s.exitCode == 0 {
				s.exitCode = 1
			}
		}
	})

	return s.exitCode, s.statusErr
}

func (s *exe) cancelFunc() error {
	s.cancel()
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}

	timeout := time.After(2 * time.Second)
	outDone := false
	errDone := false
	for !outDone || !errDone {
		select {
		case <-timeout:
			return nil
		case <-s.outDone:
			outDone = true
		case <-s.errDone:
			errDone = true
		}
	}

	return nil
}
