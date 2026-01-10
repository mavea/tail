package none

import "tail/internal/adapters/targets/console/templates/domain"

type none struct {
	name      string
	indicator domain.Indicator
}

func New(icon, name string, indicator domain.Indicator) domain.Template {
	if icon != "" {
		name = icon + " " + name
	}
	return &none{
		name:      name,
		indicator: indicator,
	}
}
func (n *none) GetHeader() string {
	return ""
}
func (n *none) GetHeaderClean() string {
	return ""
}
func (n *none) GetCellar(lines, columns, x, y uint64) string {
	return ""
}

func (n *none) GetCellarClean() string {
	return ""
}
func (n *none) FormatLine(line string) string {
	return line
}
func (n *none) StartLine() string {
	return "\033[0m"
}
func (n *none) CleanLine() string {
	return "\033[0m"
}
func (n *none) EndLine() string {
	return ""
}
