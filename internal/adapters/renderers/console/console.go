package console

import (
	"context"
	"os"
	"sync"
	"tail/internal/adapters/renderers/domain"
	"tail/internal/buffer"
	"time"
)

type renderer struct {
	ctx      context.Context
	cfg      domain.Cfg
	buffer   buffers
	error    buffers
	errorSet bool
	target   domain.Target
	print    int8
	wg       *sync.WaitGroup
	status   int
}

func New(ctx context.Context, target domain.Target, cfg domain.Cfg) (domain.Render, domain.Cancel, error) {
	result := &renderer{
		ctx:    ctx,
		cfg:    cfg,
		buffer: buffer.New(cfg),
		error:  buffer.New(cfg),
		target: target,
		wg:     &sync.WaitGroup{},

		status: -1,
	}

	result.buffer.SetDefaultStyle(result.target.GetDefaultStyle())

	if cfg.GetOutputMode() == "thread" {
		result.print = 1

		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			result.wg.Add(1)
			defer result.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					result.target.Print()
				}
			}
		}()
	}

	return result, func() {
		result.end()
	}, nil
}

func (r *renderer) Set(data string) error {
	r.buffer.Add(data)
	r.target.SetData(r.buffer.GetLast(r.cfg.GetCountLines()))

	if r.print == 0 {
		r.target.Print()
	}

	return nil
}
func (r *renderer) SetErr(err error) error {
	r.errorSet = true
	r.error.Add(err.Error())
	r.status = 1

	return nil
}

func (r *renderer) SetStatus(status int) error {
	if status == 0 && r.status != -1 {
		r.status = status
	}

	return nil
}

func (r *renderer) end() {
	r.wg.Wait()

	r.target.ClearScreen()
	if r.errorSet {
		errs := r.error.GetFull()
		if len(errs) > 0 {
			r.target.Error(r.buffer.GetFull(), errs)
		}
	}
	if r.status != 0 {
		r.target.SetStatus(r.status)
	}

	os.Exit(r.status)
}
