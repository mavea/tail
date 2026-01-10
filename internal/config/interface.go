package config

type Cfg interface {
	GetCountLines() uint64
	GetLengthLines() int
	GetSizeBuffer() uint64
	GetProcessName() string
	SetProcessName(string)
	GetProcessIcon() string
	SetProcessIcon(string)
	ReplaceTemplate(container StringValue) error
	GetTemplate() string
	ReplaceIndicator(container StringValue) error
	GetIndicator() string
	ReplaceOutputMode(container StringValue) error
	GetOutputMode() string
	IsHelp() bool
	IsVersion() bool
	GetCommand() string
	GetArgs() []string
}
