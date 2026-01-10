package console

type buffers interface {
	Add(data string)
	GetLast(len uint64) []string
	GetFull() []string
	SetDefaultStyle(string, string, string)
}
