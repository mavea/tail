package renderers

import (
	"context"
	"tail/internal/adapters/renderers/console"
	"tail/internal/adapters/renderers/domain"
)

func New(ctx context.Context, target domain.Target, cfg domain.Cfg) (domain.Render, domain.Cancel, error) {
	return console.New(ctx, target, cfg)
}
