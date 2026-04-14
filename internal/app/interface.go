package app

type Cfg interface {
	GetMaxLineCount() int
	SetMaxLineCount(int)
	GetMaxCharsPerLine() int
	SetMaxCharsPerLine(int)
	GetMaxBufferLines() uint64
	GetProcessName() string
	GetProcessIcon() string
	GetOutputTemplate() string
	GetIndicator() string
	GetOutputMode() string
	IsHelp() bool
	IsVersion() bool
	GetCommand() string
	GetArgs() []string
	Validate() error

	IsCSIEnabled() bool
	IsFullOutput() bool
}
