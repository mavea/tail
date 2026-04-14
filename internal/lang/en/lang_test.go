package en

import "testing"

func TestNewLang(t *testing.T) {
	l := NewLang()
	if l == nil {
		t.Fatalf("expected lang instance")
	}
	if l.HelpDescription.String() == "" || l.VersionDescription.String() == "" {
		t.Fatalf("expected localized messages to be filled")
	}
}
