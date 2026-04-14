package sink

import (
	"bufio"
	"context"
	formatterWindow "tail/internal/formatter/window"

	sinkConsole "tail/internal/sink/console"
	sinkDomain "tail/internal/sink/domain"
)

func NewWriter(ctx context.Context, cfg cfg, writer *bufio.Writer, window *formatterWindow.Window) (sinkDomain.Target, error) {
	return sinkConsole.New(ctx, cfg, writer, window)
}
