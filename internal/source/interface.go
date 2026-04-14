package source

type config interface {
	GetCommand() string
	GetArgs() []string
}
