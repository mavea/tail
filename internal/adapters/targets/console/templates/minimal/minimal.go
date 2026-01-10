package minimal

import (
	"tail/internal/adapters/targets/console/templates/domain"
	"tail/pkg/typo"
)

type minimal struct {
	name      string
	indicator domain.Indicator
}

func New(icon, name string, indicator domain.Indicator) domain.Template {
	if icon != "" {
		name = icon + " " + name
	}
	return &minimal{
		name:      name,
		indicator: indicator,
	}
}

func (m *minimal) GetHeader() string {
	return m.indicator.Get() + " " + m.name + typo.NewLine
}

func (m *minimal) GetHeaderClean() string {
	return m.indicator.Clean() + typo.UpAndClean
}

func (m *minimal) GetCellar(lines, columns, x, y uint64) string {
	return ""
}

func (m *minimal) GetCellarClean() string {
	return ""
}

func (m *minimal) FormatLine(line string) string {
	return line
}
func (m *minimal) StartLine() string {
	return "\033[0m"
}
func (m *minimal) CleanLine() string {
	return "\033[0m"
}
func (m *minimal) EndLine() string {
	return ""
}
