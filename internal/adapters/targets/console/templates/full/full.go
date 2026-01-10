package full

import (
	"strconv"

	"tail/internal/adapters/targets/console/templates/domain"
	"tail/pkg/typo"
)

type full struct {
	name      string
	indicator domain.Indicator
}

func New(icon, name string, indicator domain.Indicator) domain.Template {
	if icon != "" {
		name = icon + " " + name
	}
	return &full{
		name:      name,
		indicator: indicator,
	}
}

func (f *full) GetHeader() string {
	return f.indicator.Get() + " " + f.name + typo.NewLine
}

func (f *full) GetHeaderClean() string {
	return f.indicator.Clean() + typo.UpAndClean
}

func (f *full) GetCellar(lines, columns, x, y uint64) string {
	return strconv.FormatUint(lines, 10) + " : " + strconv.FormatUint(columns, 10) + " | " + strconv.FormatUint(x, 10) + " : " + strconv.FormatUint(y, 10)
}

func (f *full) GetCellarClean() string {
	return typo.UpAndClean
}

func (f *full) FormatLine(line string) string {
	return line
}
func (f *full) StartLine() string {
	return "\033[0m"
}
func (f *full) CleanLine() string {
	return "\033[0m"
}
func (f *full) EndLine() string {
	return ""
}
