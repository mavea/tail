package bootstrap

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"tail/internal/formatter"
	formatterWindow "tail/internal/formatter/window"
	langGeneral "tail/internal/lang/general"
	"tail/internal/sink"
	sinkDomain "tail/internal/sink/domain"
	"tail/internal/source"
	"tail/internal/source/general"

	"github.com/spf13/cobra"
)

var Version = "dev"
var Author = "Vitalii Makoveev (Mavea)"
var BuildTime = "unknown"

var newScanner = func(ctx context.Context, cfg Cfg) (general.Scanner, general.CanselFunc, error) {
	return source.NewScanner(ctx, cfg)
}
var newWriter = func(ctx context.Context, cfg Cfg, writer *bufio.Writer, window *formatterWindow.Window) (sinkDomain.Target, error) {
	return sink.NewWriter(ctx, cfg, writer, window)
}
var newFormatter = formatter.NewFormatter

type RunMode struct {
	cfg            Cfg
	stdin          *os.File
	stdout, stderr *bufio.Writer
	lang           *langGeneral.Lang
}

func NewRunMode(cfg Cfg, lang *langGeneral.Lang, stdin *os.File, stdout, stderr *bufio.Writer) *RunMode {
	return &RunMode{
		cfg:    cfg,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		lang:   lang,
	}
}

func (rm *RunMode) Route() func(context.Context) error {
	switch true {
	case rm.cfg.IsHelp():
		return rm.help
	case rm.cfg.IsVersion():
		return rm.version
	default:
		return rm.tail
	}
}

func (rm *RunMode) help(_ context.Context) error {
	var err error
	defer func() {
		if rm.stdout != nil {
			_ = rm.stdout.Flush()
		}
	}()

	if rm.stdout == nil {
		return fmt.Errorf("stdout is not set")
	}
	_, err = fmt.Fprint(rm.stdout, rm.lang.HelpDescription.String())
	if err != nil {
		return err
	}
	err = rm.stdout.Flush()
	if err != nil {
		return err
	}

	cmd := &cobra.Command{Use: "tail"}
	cmd.SetOut(rm.stdout)
	cmd.SetErr(rm.stdout)
	if err = cmd.Help(); err != nil {
		return err
	}
	_, err = fmt.Fprint(rm.stdout, rm.lang.HelpBottomDescription.String())
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(rm.stdout, rm.lang.HelpExample.String())

	return err
}

func (rm *RunMode) version(_ context.Context) error {
	if rm.stdout == nil {
		return fmt.Errorf("stdout is not set")
	}
	var err error
	defer func() {
		if rm.stdout != nil {
			_ = rm.stdout.Flush()
		}
	}()
	_, err = fmt.Fprintf(rm.stdout, rm.lang.VersionDescription.String(), Version, Author, BuildTime)

	return err

}

func (rm *RunMode) tail(ctx context.Context) (retErr error) {
	scanner, scannerCancel, err := newScanner(ctx, rm.cfg)
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}
	defer func() {
		_ = rm.stdout.Flush()
		cancelErr := scannerCancel()
		if cancelErr != nil {
			_, _ = fmt.Fprintf(rm.stderr, "Error: %v\n", cancelErr)
			if retErr == nil {
				retErr = cancelErr
			}
		}
	}()

	window := formatterWindow.NewWindow(rm.cfg.GetProcessIcon(), rm.cfg.GetProcessName())
	window.SetMaxSize(rm.cfg.GetMaxCharsPerLine(), rm.cfg.GetMaxLineCount())

	writer, err := newWriter(ctx, rm.cfg, rm.stdout, window)
	if err != nil {
		return fmt.Errorf("failed to create a writer: %w", err)
	}

	// Init renderer
	out, formatterCancel, err := newFormatter(ctx, writer, rm.cfg, window)
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}
	defer func() {
		cancelErr := formatterCancel()
		if cancelErr != nil {
			var exitErr formatter.ExitError
			if !errors.As(cancelErr, &exitErr) {
				_, _ = fmt.Fprintf(rm.stderr, "Error: %v\n", cancelErr)
			}
			if retErr == nil {
				retErr = cancelErr
			}
		}
	}()

	err = rm.pipe(ctx, scanner, out)
	if err != nil {
		return fmt.Errorf("failed to pipe: %w", err)
	}

	return nil
}

func (rm *RunMode) pipe(ctx context.Context, in general.Scanner, out formatter.Render) error {
	outChan := in.Out()
	errChan := in.Err()

	for outChan != nil || errChan != nil {
		select {
		case <-ctx.Done():
			if err := out.SetStatus(130); err != nil {
				return out.SetErr(err)
			}
			return nil
		case line, ok := <-outChan:
			if !ok {
				outChan = nil
				continue
			}
			if err := out.Set(line); err != nil {
				return err
			}
		case errLine, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if errLine != "" {
				if err := out.SetErrLine(errLine); err != nil {
					return err
				}
			}
		}
	}

	status, err := in.GetStatus()
	if err != nil {
		return out.SetErr(err)
	}
	if err = out.SetStatus(status); err != nil {
		return out.SetErr(err)
	}

	return nil
}
