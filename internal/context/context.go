package context

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func NewMainContext(ctx context.Context, stderr *os.File, wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	wg.Add(1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
		select {
		case s := <-quit:
			_, _ = fmt.Fprintf(stderr, "Received signal %s\n", s.String())
		case <-ctx.Done():
		}
		cancel()
		wg.Done()
	}()

	return ctx, cancel
}
