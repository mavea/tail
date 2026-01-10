package file

import (
	"bufio"
	"context"
	"os"

	"tail/internal/adapters/sources/domain"
)

type file struct {
	ctx     context.Context
	scanner *bufio.Scanner
	lineCh  chan string
	file    *os.File
}

func New(ctx context.Context, filePath string) (domain.Scanner, error) {
	// #nosec G304
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	scan := &file{
		ctx:     ctx,
		scanner: bufio.NewScanner(f),
		file:    f,
	}
	scan.scanInit()

	return scan, nil
}

func (s *file) scanInit() {
	s.lineCh = make(chan string, 1)
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
}

func (s *file) Get() (string, bool) {
	select {
	case <-s.ctx.Done():
		return "", false
	case line, ok := <-s.lineCh:
		return line, ok
	}
}

func (s *file) GetErr() (error, bool) {
	if err := s.scanner.Err(); err != nil {
		return err, true
	}

	return nil, false
}

func (s *file) GetStatus() (int, error) {
	defer func(f *os.File) {
		_ = f.Close()
	}(s.file)

	return 0, nil
}
