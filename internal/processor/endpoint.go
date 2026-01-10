package processor

import "context"

type endpoint struct {
	in  in
	out out
}

func NewService(in in, out out) Service {
	return &endpoint{
		in:  in,
		out: out,
	}
}

func (s *endpoint) Run(ctx context.Context) error {
	var (
		line    string
		ok      bool
		err     error
		errLine error
		status  int
	)
	for {
		line, ok = s.in.Get()
		if !ok && line == "" {
			break
		}

		select {
		case <-ctx.Done():
			break
		default:
			if err = s.out.Set(line); err != nil {
				return err
			}
			if !ok {
				break
			}
		}
	}
	if errLine, ok = s.in.GetErr(); ok {
		return s.out.SetErr(errLine)
	}
	if status, err = s.in.GetStatus(); err != nil {
		return s.out.SetErr(err)
	}
	if err = s.out.SetStatus(status); err != nil {
		return s.out.SetErr(err)
	}

	return nil
}
