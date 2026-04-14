package general

// Scanner defines the interface for reading data from various sources.
// Implementations should handle reading lines, errors, and status information.
type Scanner interface {
	Out() <-chan string
	Err() <-chan string
	GetStatus() (int, error)
}

type CancelFunc func() error

type CanselFunc = CancelFunc
