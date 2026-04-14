package source

import (
	"context"
	"os"

	sourceCommand "tail/internal/source/command"
	sourceFile "tail/internal/source/file"
	sourceGeneral "tail/internal/source/general"
	sourcePipe "tail/internal/source/pipe"
)

func NewScanner(ctx context.Context, cfg config) (sourceGeneral.Scanner, sourceGeneral.CanselFunc, error) {
	switch true {
	case cfg.GetCommand() != "":
		return sourceCommand.New(ctx, cfg.GetCommand())
	case len(cfg.GetArgs()) > 0:
		return sourceFile.New(ctx, cfg.GetArgs()[0])
	default:
		return sourcePipe.New(ctx, os.Stdin)
	}
}
