package formatter

import "testing"

func TestExitErrorError(t *testing.T) {
	err := ExitError{Code: 13}
	if err.Error() != "exit status 13" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
}
