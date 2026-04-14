package general

type Cfg interface {
	GetOutputTemplate() string
	GetIndicator() string
}

type Indicator interface {
	Clean(bool) string
	Get() string
}

type Template interface {
	GetHeader() string
	GetHeaderClean(bool) string
	GetCellar() string
	GetCellarClean(bool) string
	FormatLine(line string) string
	StartLine() string
	EndLine() string
	CleanLine() string
}

type Window interface {
	GetPosition() (uint64, uint64)
	GetBufferSize() (uint64, uint64)
	GetIcon() string
	GetTitle() string
}
