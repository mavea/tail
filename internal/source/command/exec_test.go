package command_test

import (
	"context"
	"runtime"
	"sort"
	"strings"
	"tail/internal/source/command"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func shellCommand(script string) (string, []string) {
	if runtime.GOOS == "windows" {
		return "cmd", []string{"/C", script}
	}

	return "sh", []string{"-c", script}
}

var _ = Describe("Exec", func() {
	Context("Multiple GetStatus calls", func() {
		It("should return the same result for multiple calls", func() {
			script := "echo test; exit 0"
			if runtime.GOOS == "windows" {
				script = "echo test & exit /b 0"
			}
			cmd, args := shellCommand(script)
			s, _, err := command.New(context.Background(), cmd, args...)
			Expect(err).NotTo(HaveOccurred())

			for range s.Out() {
			}
			for range s.Err() {
			}

			status1, err1 := s.GetStatus()
			Expect(err1).NotTo(HaveOccurred())
			Expect(status1).To(BeZero())

			status2, err2 := s.GetStatus()
			if err1 == nil {
				Expect(err2).NotTo(HaveOccurred())
			} else {
				Expect(err2).To(MatchError(err1))
			}
			Expect(status2).To(Equal(status1))
		})
	})

	Context("cancelFunc and deadlocks", func() {
		It("should not deadlock in cancelFunc if the process is stuck", func() {
			longRun := "sleep 10"
			if runtime.GOOS == "windows" {
				longRun = "ping -n 10 127.0.0.1 > nul"
			}
			cmd, args := shellCommand(longRun)
			_, cancel, err := command.New(context.Background(), cmd, args...)
			Expect(err).NotTo(HaveOccurred())

			done := make(chan bool)
			go func() {
				_ = cancel()
				done <- true
			}()

			select {
			case <-done:
			case <-time.After(5 * time.Second):
				Fail("cancelFunc deadlocked")
			}
		})
	})

	Context("Out/Err channels", func() {
		It("should fail fast for unresolved simple command names", func() {
			s, cancel, err := command.New(context.Background(), "tail_command_that_does_not_exist_123456")
			Expect(s).To(BeNil())
			Expect(cancel).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to resolve executable"))
		})

		It("should execute a single command string", func() {
			s, cancel, err := command.New(context.Background(), "echo single-string")
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = cancel()
			}()

			var out []string
			for line := range s.Out() {
				out = append(out, line)
			}
			for range s.Err() {
			}

			Expect(out).To(ContainElement("single-string"))

			status, statusErr := s.GetStatus()
			Expect(statusErr).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should route stdout to Out and stderr to Err", func() {
			script := "echo out1; echo err1 1>&2; echo out2; echo err2 1>&2"
			if runtime.GOOS == "windows" {
				script = "(echo out1 & echo err1 1>&2 & echo out2 & echo err2 1>&2)"
			}
			cmd, args := shellCommand(script)
			s, cancel, err := command.New(context.Background(), cmd, args...)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = cancel()
			}()

			var out []string
			for line := range s.Out() {
				out = append(out, strings.TrimSpace(line))
			}

			var stderr []string
			for line := range s.Err() {
				stderr = append(stderr, strings.TrimSpace(line))
			}

			sort.Strings(out)
			sort.Strings(stderr)
			Expect(out).To(Equal([]string{"out1", "out2"}))
			Expect(stderr).To(Equal([]string{"err1", "err2"}))

			status, statusErr := s.GetStatus()
			Expect(statusErr).NotTo(HaveOccurred())
			Expect(status).To(Equal(0))
		})

		It("should return non-zero status for a failed command", func() {
			script := "echo boom 1>&2; exit 7"
			if runtime.GOOS == "windows" {
				script = "echo boom 1>&2 & exit /b 7"
			}
			cmd, args := shellCommand(script)
			s, cancel, err := command.New(context.Background(), cmd, args...)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = cancel()
			}()

			for range s.Out() {
			}
			for range s.Err() {
			}

			status, statusErr := s.GetStatus()
			Expect(statusErr).NotTo(HaveOccurred())
			Expect(status).To(Equal(7))
		})
	})
})
