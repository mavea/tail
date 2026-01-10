package pipe

import (
	"bufio"
	"context"
	"os"

	"tail/internal/adapters/sources/domain"
)

type pipe struct {
	ctx     context.Context
	scanner *bufio.Scanner
	lineCh  chan string
}

func New(ctx context.Context, in *os.File) (domain.Scanner, error) {
	scan := &pipe{
		ctx:     ctx,
		scanner: bufio.NewScanner(in),
	}
	scan.scanInit()

	return scan, nil
}

func (s *pipe) scanInit() {
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

func (s *pipe) Get() (string, bool) {
	select {
	case <-s.ctx.Done():
		return "", false
	case line, ok := <-s.lineCh:
		return line, ok
	}
}

func (s *pipe) GetErr() (error, bool) {
	return nil, false
}

func (s *pipe) GetStatus() (int, error) {
	/*	cmd := exec.Command("bash", "-c", "echo $?")
		output, err := cmd.Output()
		if err != nil {
			return 2, err
		}

		code, err := strconv.Atoi(strings.TrimSpace(string(output)))
		if err != nil {
			return 3, err
		}
	*/
	return 0, nil
}
