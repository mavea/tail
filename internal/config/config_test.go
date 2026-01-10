package config

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockStringValue struct {
	value string
	err   error
}

func (m *mockStringValue) String() string {
	return m.value
}

func (m *mockStringValue) Set(v string) error {
	if m.err != nil {
		return m.err
	}
	m.value = v
	return nil
}

func (m *mockStringValue) Type() string {
	return "mockStringValue"
}

var _ = Describe("Config", func() {
	Context("GetCountLines", func() {
		DescribeTable("Sunny",
			func(initial, expected uint64) {
				cfg := &config{countLines: initial}
				Expect(cfg.GetCountLines()).To(Equal(expected))
			},
			Entry("zero", uint64(0), uint64(0)),
			Entry("not zero", uint64(10), uint64(10)),
		)
	})
	Context("GetLengthLines", func() {
		DescribeTable("Sunny",
			func(initial int, expected int) {
				cfg := &config{lengthLines: initial}
				Expect(cfg.GetLengthLines()).To(Equal(expected))
			},
			Entry("zero", int(0), int(0)),
			Entry("not zero", int(50), int(50)),
		)
	})
	Context("GetSizeBuffer", func() {
		DescribeTable("Sunny",
			func(initial, expected uint64) {
				cfg := &config{sizeBuffer: initial}
				Expect(cfg.GetSizeBuffer()).To(Equal(expected))
			},
			Entry("zero", uint64(0), uint64(0)),
			Entry("not zero", uint64(100), uint64(100)),
		)
	})
	Context("GetProcessName", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{processName: initial}
				Expect(cfg.GetProcessName()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "testProcess", "testProcess"),
		)
	})
	Context("SetProcessName", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{}
				cfg.SetProcessName(initial)
				Expect(cfg.GetProcessName()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "newProcess", "newProcess"),
		)
	})
	Context("GetProcessIcon", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{processIcon: initial}
				Expect(cfg.GetProcessIcon()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "✅", "✅"),
		)
	})
	Context("SetProcessIcon", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{}
				cfg.SetProcessIcon(initial)
				Expect(cfg.GetProcessIcon()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "🔨", "🔨"),
		)
	})
	Context("IsHelp", func() {
		DescribeTable("Sunny",
			func(initial, expected bool) {
				cfg := &config{help: initial}
				Expect(cfg.IsHelp()).To(Equal(expected))
			},
			Entry("false", false, false),
			Entry("true", true, true),
		)
	})
	Context("IsVersion", func() {
		DescribeTable("Sunny",
			func(initial, expected bool) {
				cfg := &config{version: initial}
				Expect(cfg.IsVersion()).To(Equal(expected))
			},
			Entry("false", false, false),
			Entry("true", true, true),
		)
	})
	Context("GetCommand", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{command: initial}
				Expect(cfg.GetCommand()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "testCommand", "testCommand"),
		)
	})
	Context("GetArgs", func() {
		DescribeTable("Sunny",
			func(initial, expected []string) {
				cfg := &config{args: initial}
				Expect(cfg.GetArgs()).To(Equal(expected))
			},
			Entry("nil", nil, nil),
			Entry("empty", []string{}, []string{}),
			Entry("one", []string{"arg"}, []string{"arg"}),
			Entry("two", []string{"arg1", "arg2"}, []string{"arg1", "arg2"}),
		)
	})
	Context("ReplaceOutputMode", func() {
		DescribeTable("Sunny",
			func(initial StringValue, expected string, wantErr error) {
				cfg := &config{outputMode: &mockStringValue{value: "oldOutput"}}

				err := cfg.ReplaceOutputMode(initial)
				if wantErr != nil {
					Expect(err).To(Equal(wantErr))
				} else {
					Expect(err).To(BeNil())
				}
				if err == nil {
					Expect(cfg.GetOutputMode()).To(Equal(expected))
				}
			},
			Entry("valid", &mockStringValue{value: "newOutput"}, "oldOutput", nil),
			Entry("error", &mockStringValue{err: errors.New("failure")}, "", errors.New("failure")),
		)
	})
	Context("GetOutputMode", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{outputMode: &mockStringValue{value: initial}}
				Expect(cfg.GetOutputMode()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "mode", "mode"),
		)
	})
	Context("ReplaceTemplate", func() {
		DescribeTable("Sunny",
			func(initial StringValue, expected string, wantErr error) {
				cfg := &config{template: &mockStringValue{value: "oldTemplate"}}

				err := cfg.ReplaceTemplate(initial)
				if wantErr == nil {
					Expect(err).To(BeNil())
				} else {
					Expect(err).To(Equal(wantErr))
				}
				if err == nil {
					Expect(cfg.GetTemplate()).To(Equal(expected))
				}
			},
			Entry("valid", &mockStringValue{value: "newTemplate"}, "oldTemplate", nil),
			Entry("error", &mockStringValue{err: errors.New("failure")}, "", errors.New("failure")),
		)
	})
	Context("GetTemplate", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{template: &mockStringValue{value: initial}}
				Expect(cfg.GetTemplate()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "template", "template"),
		)
	})
	Context("ReplaceIndicator", func() {
		DescribeTable("Sunny",
			func(initial StringValue, expected string, wantErr error) {
				cfg := &config{indicator: &mockStringValue{value: "oldTemplate"}}

				err := cfg.ReplaceIndicator(initial)
				if wantErr != nil {
					Expect(err).To(Equal(wantErr))
				} else {
					Expect(err).To(BeNil())
				}
				if err == nil {
					Expect(cfg.GetIndicator()).To(Equal(expected))
				}
			},
			Entry("valid", &mockStringValue{value: "newIndicator"}, "oldTemplate", nil),
			Entry("error", &mockStringValue{err: errors.New("failure")}, "", errors.New("failure")),
		)
	})
	Context("GetIndicator", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				cfg := &config{indicator: &mockStringValue{value: initial}}
				Expect(cfg.GetIndicator()).To(Equal(expected))
			},
			Entry("zero", "", ""),
			Entry("not zero", "mode", "mode"),
		)
	})
})
