package sources

type config interface {
	GetCommand() string
	GetArgs() []string
}
