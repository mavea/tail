package pipe

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	sourceGeneral "tail/internal/source/general"
)

const chanCache = 1
const maxScanTokenSize = 1024 * 1024
const closePipeTimeout = 200 * time.Millisecond

type pipe struct {
	ctx    context.Context
	cancel context.CancelFunc

	scanner   *bufio.Scanner
	in        *os.File
	outChan   chan string
	errChan   chan string
	statusErr error
	mu        sync.Mutex

	scanDone chan struct{}
}

func New(ctx context.Context, in *os.File) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
	s := &pipe{
		scanner:  bufio.NewScanner(in),
		in:       in,
		outChan:  make(chan string, chanCache),
		errChan:  make(chan string, chanCache),
		scanDone: make(chan struct{}),
	}
	s.scanner.Buffer(make([]byte, 64*1024), maxScanTokenSize)
	s.ctx, s.cancel = context.WithCancel(ctx)

	s.scan()

	return s, s.cancelFunc, nil
}

func (s *pipe) scan() {
	go func() {
		defer func() {
			close(s.outChan)
			close(s.errChan)
			close(s.scanDone)
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
			s.statusErr = fmt.Errorf("scanner error: %w", err)
			s.mu.Unlock()
		}
	}()
}

func (s *pipe) Out() <-chan string {
	return s.outChan
}

func (s *pipe) Err() <-chan string {
	return s.errChan
}

func (s *pipe) GetStatus() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.statusErr != nil {
		return 1, s.statusErr
	}
	return 0, nil
}

func (s *pipe) waitGroup() {
	select {
	case <-s.scanDone:
	case <-time.After(2 * time.Second):
	}
}

func (s *pipe) cancelFunc() error {
	s.cancel()
	if s.in != nil && s.in != os.Stdin {
		closeErr := make(chan error, 1)
		go func() {
			closeErr <- s.in.Close()
		}()

		select {
		case err := <-closeErr:
			if err != nil {
				return fmt.Errorf("failed to close pipe: %w", err)
			}
		case <-time.After(closePipeTimeout):
			// Avoid deadlock when pipe close blocks while scanner is in a blocking read.
		}
	}
	s.waitGroup()

	return nil
}
