package domain

type Target interface {
	GetDefaultStyle() (string, string, string)
	Print()
	SetData(data []string)
	ClearScreen()
	Error(buffer []string, err []string)
	SetStatus(status int)
}
