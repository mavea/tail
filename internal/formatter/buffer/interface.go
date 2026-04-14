package buffer

type config interface {
	GetMaxLineCount() int
	GetMaxCharsPerLine() int
	GetMaxBufferLines() uint64

	GetProcessName() string
	GetProcessIcon() string

	IsCSIEnabled() bool
	IsFullOutput() bool
}

type Buffer interface {
	Add(line string)
	GetLast(len int) []string
	GetFull() []string
	SetDefaultStyle(string, string, string)
}

type Window interface {
	SetPosition(x int, y uint64)
	SetBufferSize(countBufferLines uint64, countBufferColumns int)
	SetIcon(icon string)
	SetTitle(title string)
	Height() uint64
}
