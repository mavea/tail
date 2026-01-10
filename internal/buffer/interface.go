package buffer

type config interface {
	GetCountLines() uint64
	GetLengthLines() int
	GetSizeBuffer() uint64

	GetProcessName() string
	SetProcessName(string)
	GetProcessIcon() string
	SetProcessIcon(string)
}

type Buffer interface {
	Add(line string)
	GetLast(len uint64) []string
	GetFull() []string
	/*Err(data error)
	GetErr() []string*/
	SetDefaultStyle(string, string, string)
}
