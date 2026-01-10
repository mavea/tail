package style

type Style interface {
	Set([]uint64) Style
	String() string
}
