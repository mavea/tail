package targets

import (
	"bufio"
	"tail/internal/adapters/targets/console"
	"tail/internal/adapters/targets/domain"
)

func New(cfg cfg, writer *bufio.Writer) (domain.Target, error) {
	return console.New(cfg, writer)
}
