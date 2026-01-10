package templates

import (
	"tail/internal/adapters/targets/console/indicator"
	"tail/internal/adapters/targets/console/templates/domain"
	"tail/internal/adapters/targets/console/templates/full"
	"tail/internal/adapters/targets/console/templates/minimal"
	"tail/internal/adapters/targets/console/templates/none"
)

func NewTemplate(cfg domain.Cfg, indicate indicator.Indicator) (domain.Template, error) {
	switch cfg.GetTemplate() {
	case "minimal":
		return minimal.New(cfg.GetProcessIcon(), cfg.GetProcessName(), indicate), nil
	case "full":
		return full.New(cfg.GetProcessIcon(), cfg.GetProcessName(), indicate), nil
	default:
		return none.New(cfg.GetProcessIcon(), cfg.GetProcessName(), indicate), nil
	}
}
