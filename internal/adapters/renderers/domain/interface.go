package domain

type Render interface {
	Set(data string) error
	SetErr(err error) error
	SetStatus(status int) error
}

type Cancel func()

type Target interface {
	GetDefaultStyle() (string, string, string)
	Print()
	SetData(data []string)
	ClearScreen()
	Error(buffer []string, err []string)
	SetStatus(status int)
}

type Cfg interface {
	GetCountLines() uint64
	GetLengthLines() int
	GetSizeBuffer() uint64
	GetProcessName() string
	SetProcessName(string)
	GetProcessIcon() string
	SetProcessIcon(string)
	GetTemplate() string
	GetIndicator() string
	GetOutputMode() string
}
