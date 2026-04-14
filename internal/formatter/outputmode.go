package formatter

import (
	"errors"
	"strings"

	configGeneral "tail/internal/config/general"
)

const (
	outputModeDirect = "direct"
	outputModeThread = "thread"
)

type outputMode struct {
	value string
}

func NewOutputMode() configGeneral.StringValue {
	return &outputMode{value: outputModeDirect}
}

func (m *outputMode) String() string {
	return m.value
}

func (m *outputMode) Set(s string) error {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", outputModeDirect:
		m.value = outputModeDirect
	case outputModeThread:
		m.value = outputModeThread
	default:
		return errors.New("invalid output mode value: " + s)
	}

	return nil
}

func (m *outputMode) Type() string {
	return "string"
}

func (m *outputMode) Validate() bool {
	switch m.value {
	case outputModeDirect, outputModeThread:
		return true
	default:
		return false
	}
}
