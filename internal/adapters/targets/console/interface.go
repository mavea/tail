package console

import "tail/internal/config"

type cfg interface {
	GetCountLines() uint64
	GetLengthLines() int
	GetProcessName() string
	GetProcessIcon() string
	ReplaceTemplate(container config.StringValue) error
	GetTemplate() string
	ReplaceIndicator(container config.StringValue) error
	GetIndicator() string
}
