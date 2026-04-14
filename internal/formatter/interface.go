package formatter

// Render defines the interface for processing and rendering output data.
// Implementations handle setting lines, errors, and status information.
type Render interface {
	// Set processes a line of output data.
	Set(data string) error
	// SetErrLine processes a stderr line without forcing a failed exit status.
	SetErrLine(data string) error
	// SetErr handles error information from the source.
	SetErr(err error) error
	// SetStatus handles status information from the source.
	SetStatus(status int) error
}

// Cancel is a function type for cancelling rendering operations.
type Cancel func() error

// Out defines the interface for output destinations.
// Implementations handle displaying data to various output targets.
type Out interface {
	// GetDefaultStyle returns the default styling strings (start, clean, end).
	GetDefaultStyle() (string, string, string)
	// Print renders the current buffer to the output target.
	Print() error
	// SetData updates the data to be displayed.
	SetData(data []string)
	// ClearScreen clears the output display.
	ClearScreen() error
	// Error displays error information.
	Error(buffer []string, err []string) error
	// PrintFull displays full information.
	PrintFull(buffer []string) error
	// SetStatus displays status information.
	SetStatus(status int) error
}

// Cfg defines the interface for configuration access.
type Cfg interface {
	GetMaxLineCount() int
	GetMaxCharsPerLine() int
	GetMaxBufferLines() uint64
	GetProcessName() string
	GetProcessIcon() string
	GetOutputTemplate() string
	GetIndicator() string
	GetOutputMode() string

	IsCSIEnabled() bool
	IsFullOutput() bool
}

type buffers interface {
	Add(data string)
	GetLast(len int) []string
	GetFull() []string
	SetDefaultStyle(string, string, string)
}
