package minimal

import (
	sinkConsoleTemplateGeneral "tail/internal/sink/console/template/general"
	typo "tail/internal/typo"
)

type minimal struct {
	indicator sinkConsoleTemplateGeneral.Indicator
	window    sinkConsoleTemplateGeneral.Window
}

func New(indicator sinkConsoleTemplateGeneral.Indicator, window sinkConsoleTemplateGeneral.Window) sinkConsoleTemplateGeneral.Template {
	return &minimal{
		indicator: indicator,
		window:    window,
	}
}

func (m *minimal) GetHeader() string {
	return typo.Zero + m.indicator.Get() + " " + m.window.GetIcon() + m.window.GetTitle() + typo.NewLine
}

func (m *minimal) GetHeaderClean(firstLine bool) string {
	if firstLine {
		return ""
	}

	return m.indicator.Clean(firstLine) + typo.UpAndClean
}

func (m *minimal) GetCellar() string {
	return ""
}

func (m *minimal) GetCellarClean(_ bool) string {
	return ""
}

func (m *minimal) FormatLine(line string) string {
	return line
}

func (m *minimal) StartLine() string {
	return typo.Zero
}

func (m *minimal) CleanLine() string {
	return ""
}

func (m *minimal) EndLine() string {
	return ""
}
