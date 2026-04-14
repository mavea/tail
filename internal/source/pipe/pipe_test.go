package pipe_test

import (
	"context"
	"os"
	"time"

	sourcePipe "tail/internal/source/pipe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipe source", func() {
	It("reads all lines from pipe throughout and returns zero status", func() {
		r, w, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = r.Close()
		}()

		s, cancel, err := sourcePipe.New(context.Background(), r)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = cancel()
		}()

		go func() {
			_, _ = w.WriteString("alpha\nbeta\ngamma\n")
			_ = w.Close()
		}()

		var out []string
		for line := range s.Out() {
			out = append(out, line)
		}
		Expect(out).To(Equal([]string{"alpha", "beta", "gamma"}))

		select {
		case _, ok := <-s.Err():
			Expect(ok).To(BeFalse())
		case <-time.After(200 * time.Millisecond):
			Fail("Err channel must be closed for pipe source")
		}

		status, statusErr := s.GetStatus()
		Expect(statusErr).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
	})

	It("cancel does not deadlock", func() {
		r, w, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = r.Close()
			_ = w.Close()
		}()

		s, cancel, err := sourcePipe.New(context.Background(), r)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).NotTo(BeNil())

		done := make(chan struct{})
		go func() {
			_ = cancel()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(3 * time.Second):
			Fail("cancel deadlocked")
		}
	})
})
