package domain

// Target defines the interface for output destinations.
// Implementations handle displaying data to various output targets.
type Target interface {
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
