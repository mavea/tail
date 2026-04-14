package none

import (
	sinkConsoleGeneral "tail/internal/sink/console/template/general"
	"tail/internal/typo"
)

type none struct {
	indicator sinkConsoleGeneral.Indicator
}

func New(indicator sinkConsoleGeneral.Indicator, _ sinkConsoleGeneral.Window) sinkConsoleGeneral.Template {
	return &none{
		indicator: indicator,
	}
}

func (n *none) GetHeader() string {
	return ""
}

func (n *none) GetHeaderClean(_ bool) string {
	return ""
}

func (n *none) GetCellar() string {
	return ""
}

func (n *none) GetCellarClean(_ bool) string {
	return ""
}

func (n *none) FormatLine(line string) string {
	return line
}

func (n *none) StartLine() string {
	return typo.Zero
}

func (n *none) CleanLine() string {
	return ""
}

func (n *none) EndLine() string {
	return ""
}
