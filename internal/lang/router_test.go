package lang

import (
	"testing"

	langEn "tail/internal/lang/en"
	langRu "tail/internal/lang/ru"
)

func TestNewLangPackage(t *testing.T) {
	enPack, err := NewLangPackage(langEn.Code)
	if err != nil || enPack == nil {
		t.Fatalf("expected en package, got err=%v", err)
	}

	ruPack, err := NewLangPackage(langRu.Code)
	if err != nil || ruPack == nil {
		t.Fatalf("expected ru package, got err=%v", err)
	}

	if _, err = NewLangPackage("de"); err == nil {
		t.Fatalf("expected unsupported language error")
	}
}

func TestGetLang(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "en", want: "en"},
		{in: "EN_us.UTF-8", want: "en"},
		{in: "ru_RU.UTF-8", want: "ru"},
		{in: "", want: string(Default)},
		{in: "xx", want: string(Default)},
	}

	for _, tc := range tests {
		got := string(GetLang(tc.in))
		if got != tc.want {
			t.Fatalf("GetLang(%q)=%q, want %q", tc.in, got, tc.want)
		}
	}
}
