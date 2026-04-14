package template

import (
	sinkConsoleIndicator "tail/internal/sink/console/indicator"
	sinkConsoleTemplateFull "tail/internal/sink/console/template/full"
	sinkConsoleTemplateGeneral "tail/internal/sink/console/template/general"
	sinkConsoleTemplateMinimal "tail/internal/sink/console/template/minimal"
	sinkConsoleTemplateNone "tail/internal/sink/console/template/none"
)

func NewTemplate(cfg sinkConsoleTemplateGeneral.Cfg, indicate sinkConsoleIndicator.Indicator, window sinkConsoleTemplateGeneral.Window) (sinkConsoleTemplateGeneral.Template, error) {
	switch cfg.GetOutputTemplate() {
	case "full":
		return sinkConsoleTemplateFull.New(indicate, window), nil
	case "minimal":
		return sinkConsoleTemplateMinimal.New(indicate, window), nil
	default:
		return sinkConsoleTemplateNone.New(indicate, window), nil
	}
}
