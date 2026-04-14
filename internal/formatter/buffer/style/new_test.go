package style

import "testing"

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New should return style")
	}
	s2 := s.Set([]uint64{31})
	if s2 == nil {
		t.Fatal("Set should return style")
	}
	if s2.String() == s.String() {
		t.Fatal("style string should change after Set")
	}
}
