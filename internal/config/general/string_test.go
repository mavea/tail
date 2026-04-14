package general

import "testing"

func TestDefaultStringImplementation(t *testing.T) {
	ds := &DefaultString{}
	if ds.String() != "" {
		t.Fatalf("unexpected empty string value")
	}

	if ds.Set("hello") != nil {
		t.Fatal("set should not fail")
	}
	if ds.String() != "hello" {
		t.Fatalf("unexpected string value: %s", ds.String())
	}

	if ds.Type() != "string" {
		t.Fatalf("unexpected type: %s", ds.Type())
	}

	if !ds.Validate() {
		t.Fatal("validate should return true")
	}
}
