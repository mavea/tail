package file

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	sourceGeneral "tail/internal/source/general"
)

const chanCache = 1
const maxScanTokenSize = 1024 * 1024

type file struct {
	ctx    context.Context
	cancel context.CancelFunc

	scanner *bufio.Scanner
	in      *os.File
	outChan chan string
	errChan chan string

	exitCode  int
	statusErr error
	mu        sync.Mutex

	wg *sync.WaitGroup
}

func New(ctx context.Context, filePath string) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
	// #nosec G304 -- File mode intentionally opens user-provided path from CLI arguments.
	openFile, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	s := &file{
		scanner:  bufio.NewScanner(openFile),
		in:       openFile,
		outChan:  make(chan string, chanCache),
		errChan:  make(chan string, chanCache),
		wg:       &sync.WaitGroup{},
		exitCode: 0,
	}
	s.scanner.Buffer(make([]byte, 64*1024), maxScanTokenSize)
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.scan()

	return s, s.cancelFunc, nil
}

func (s *file) scan() {
	s.wg.Add(1)
	go func() {
		defer func() {
			close(s.outChan)
			close(s.errChan)
			s.wg.Done()
		}()

		for s.scanner.Scan() {
			select {
			case <-s.ctx.Done():
				return
			case s.outChan <- s.scanner.Text():
			}
		}

		if err := s.scanner.Err(); err != nil {
			s.mu.Lock()
			s.exitCode = 1
			s.statusErr = fmt.Errorf("scanner error: %w", err)
			s.mu.Unlock()
			select {
			case s.errChan <- "scanner error: " + err.Error():
			default:
			}
		}
	}()
}

func (s *file) Out() <-chan string {
	return s.outChan
}

func (s *file) Err() <-chan string {
	return s.errChan
}

func (s *file) GetStatus() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.exitCode, s.statusErr
}

func (s *file) cancelFunc() error {
	err := s.in.Close()
	if err != nil {
		s.mu.Lock()
		s.exitCode = 1
		s.statusErr = fmt.Errorf("failed to close file: %w", err)
		s.mu.Unlock()

		return fmt.Errorf("failed to close file: %w", err)
	}
	s.cancel()
	s.wg.Wait()

	return nil
}
