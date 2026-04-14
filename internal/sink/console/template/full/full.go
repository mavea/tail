package full

import (
	"strconv"

	sinkConsoleTemplateGeneral "tail/internal/sink/console/template/general"
	typo "tail/internal/typo"
)

type full struct {
	window    sinkConsoleTemplateGeneral.Window
	indicator sinkConsoleTemplateGeneral.Indicator
}

func New(indicator sinkConsoleTemplateGeneral.Indicator, window sinkConsoleTemplateGeneral.Window) sinkConsoleTemplateGeneral.Template {
	return &full{
		indicator: indicator,
		window:    window,
	}
}

func (f *full) GetHeader() string {
	return typo.Zero + f.indicator.Get() + " " + f.window.GetIcon() + f.window.GetTitle() + typo.NewLine
}

func (f *full) GetHeaderClean(firstLine bool) string {
	if firstLine {
		return ""
	}

	return f.indicator.Clean(firstLine) + typo.UpAndClean
}

func (f *full) GetCellar() string {
	lines, columns := f.window.GetBufferSize()
	x, y := f.window.GetPosition()
	return typo.Zero + strconv.FormatUint(lines, 10) + " : " + strconv.FormatUint(columns, 10) + " | " + strconv.FormatUint(x, 10) + " : " + strconv.FormatUint(y, 10) + typo.NewLine
}

func (f *full) GetCellarClean(firstLine bool) string {
	if firstLine {
		return ""
	}

	return typo.UpAndClean
}

func (f *full) FormatLine(line string) string {
	return line
}

func (f *full) StartLine() string {
	return typo.Zero
}

func (f *full) CleanLine() string {
	return ""
}

func (f *full) EndLine() string {
	return ""
}
