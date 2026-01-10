package sources

import (
	"context"
	"os"

	"tail/internal/adapters/sources/command"
	"tail/internal/adapters/sources/domain"
	"tail/internal/adapters/sources/file"
	"tail/internal/adapters/sources/pipe"
)

func NewScanner(ctx context.Context, cfg config) (domain.Scanner, error) {
	switch true {
	case cfg.GetCommand() != "":
		return command.New(ctx, cfg.GetCommand())
	case len(cfg.GetArgs()) > 0:
		return file.New(ctx, cfg.GetArgs()[0])
	default:
		return pipe.New(ctx, os.Stdin)
	}
}
