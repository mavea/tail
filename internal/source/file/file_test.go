package file_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	sourceFile "tail/internal/source/file"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File source", func() {
	It("reads all lines from the file throughout and returns zero status", func() {
		tmpDir := GinkgoT().TempDir()
		filePath := filepath.Join(tmpDir, "input.log")
		err := os.WriteFile(filePath, []byte("one\ntwo\nthree\n"), 0o600)
		Expect(err).NotTo(HaveOccurred())

		s, cancel, err := sourceFile.New(context.Background(), filePath)
		Expect(err).NotTo(HaveOccurred())

		var out []string
		for line := range s.Out() {
			out = append(out, line)
		}
		Expect(out).To(Equal([]string{"one", "two", "three"}))

		Expect(cancel()).To(Succeed())

		select {
		case _, ok := <-s.Err():
			Expect(ok).To(BeFalse())
		case <-time.After(2 * time.Second):
			Fail("Err channel must be closed for file source")
		}

		status, statusErr := s.GetStatus()
		Expect(statusErr).NotTo(HaveOccurred())
		Expect(status).To(Equal(0))
	})

	It("returns an error when file does not exist", func() {
		s, cancel, err := sourceFile.New(context.Background(), filepath.Join("not", "exists", "file.log"))
		Expect(s).To(BeNil())
		Expect(cancel).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(strings.Contains(err.Error(), "failed to open file")).To(BeTrue())
	})
})
