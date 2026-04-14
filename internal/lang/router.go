package lang

import (
	"fmt"
	"strings"

	langEn "tail/internal/lang/en"
	langGeneral "tail/internal/lang/general"
	langRu "tail/internal/lang/ru"
)

const Default = langEn.Code

func NewLangPackage(lang langGeneral.Code) (*langGeneral.Lang, error) {
	switch lang {
	case langEn.Code:
		return langEn.NewLang(), nil
	case langRu.Code:
		return langRu.NewLang(), nil
	default:
		return nil, fmt.Errorf("unsupported language code: %s", lang)
	}
}

func GetLang(lang string) langGeneral.Code {
	langShort := lang
	if lang != "" && len(lang) > 2 {
		langShort = lang[:2]
	}

	switch strings.ToLower(langShort) {
	case string(langEn.Code):
		return langEn.Code
	case string(langRu.Code):
		return langRu.Code
	default:
		return Default
	}
}
