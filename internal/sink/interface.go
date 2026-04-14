package sink

type cfg interface {
	GetMaxLineCount() int
	GetMaxCharsPerLine() int
	GetProcessName() string
	GetProcessIcon() string
	GetOutputTemplate() string
	GetIndicator() string
}

type Window interface {
	GetPosition() (uint64, uint64)
	GetBufferSize() (uint64, uint64)
	GetIcon() string
	GetTitle() string
}
