package console

type cfg interface {
	GetMaxLineCount() int
	GetMaxCharsPerLine() int
	GetOutputTemplate() string
	GetIndicator() string
}
