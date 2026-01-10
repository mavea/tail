package command

import (
	"bufio"
	"context"
	"errors"
	"os/exec"

	"tail/internal/adapters/sources/domain"
)

type exe struct {
	ctx        context.Context
	scanner    *bufio.Scanner
	lineCh     chan string
	scannerErr *bufio.Scanner
	lineErrCh  chan string
	cmd        *exec.Cmd
}

func New(ctx context.Context, programPath string) (domain.Scanner, error) {
	cmd := exec.Command(programPath)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	scan := &exe{
		ctx:        ctx,
		scanner:    bufio.NewScanner(stdoutPipe),
		scannerErr: bufio.NewScanner(stderrPipe),
		cmd:        cmd,
	}
	scan.scanInit()

	return scan, nil
}

func (s *exe) scanInit() {
	s.lineCh = make(chan string, 1)
	s.lineErrCh = make(chan string, 1)
	go func() {
		defer close(s.lineCh)
		for s.scanner.Scan() {
			select {
			case <-s.ctx.Done():
				return
			case s.lineCh <- s.scanner.Text():
			}
		}
	}()
	go func() {
		defer close(s.lineErrCh)
		for s.scannerErr.Scan() {
			select {
			case <-s.ctx.Done():
				return
			case s.lineErrCh <- s.scannerErr.Text():
			}
		}
	}()
}

func (s *exe) Get() (string, bool) {
	select {
	case <-s.ctx.Done():
		return "", false
	case line, ok := <-s.lineCh:
		return line, ok
	}
}

func (s *exe) GetErr() (error, bool) {
	select {
	case <-s.ctx.Done():
		return nil, false
	case line, ok := <-s.lineErrCh:
		if line != "" {
			return errors.New(line), ok
		}
		return nil, ok
	}
}

func (s *exe) GetStatus() (int, error) {
	err := s.cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return 1, nil
	}
	return 0, nil
}
