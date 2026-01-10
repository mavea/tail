package domain

import "tail/internal/config"

type Cfg interface {
	GetProcessName() string
	GetProcessIcon() string
	ReplaceTemplate(container config.StringValue) error
	GetTemplate() string
	ReplaceIndicator(container config.StringValue) error
	GetIndicator() string
}

type Indicator interface {
	Clean() string
	Get() string
}

type Template interface {
	GetHeader() string
	GetHeaderClean() string
	GetCellar(lines, columns, x, y uint64) string
	GetCellarClean() string
	FormatLine(line string) string
	StartLine() string
	EndLine() string
	CleanLine() string
}
