package general

import "testing"

func TestMessageString(t *testing.T) {
	m := Message("hello")
	if m.String() != "hello" {
		t.Fatalf("unexpected message string: %s", m.String())
	}
}
