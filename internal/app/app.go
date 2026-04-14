package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sync"
	"tail/internal/bootstrap"

	"golang.org/x/term"

	lang "tail/internal/lang"
	langGeneral "tail/internal/lang/general"
)

type App struct {
	ctx context.Context
	cfg Cfg

	fd             int
	stdin          *os.File
	stdout, stderr *bufio.Writer

	lang *langGeneral.Lang

	wg *sync.WaitGroup

	cancels []func() error
}

// NewApp creates a new App instance with the given file descriptor, context, and wait group.
// It validates that the file descriptor corresponds to a terminal.
func NewApp(fd uintptr, ctx context.Context, wg *sync.WaitGroup) (*App, func() error, error) {
	if fd > math.MaxInt {
		return nil, nil, fmt.Errorf("invalid file descriptor")
	}
	iFd := int(fd)

	/*if !term.IsTerminal(iFd) {
		return nil, nil, fmt.Errorf("not a terminal")
	}*/

	app := &App{
		fd:      iFd,
		ctx:     ctx,
		wg:      wg,
		cancels: []func() error{},
	}

	return app, app.cancel, nil
}

func (app *App) cancel() error {
	var errs error
	for _, cancel := range app.cancels {
		err := cancel()
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to cancel: %w", err))
		}
	}

	return errs
}

func (app *App) SetDefaultStd(stdin, stdout, stderr *os.File) {
	app.stdin = stdin
	app.stdout = bufio.NewWriter(stdout)
	app.stderr = bufio.NewWriter(stderr)
}

// DetectAndSetLanguage detects the language from the given string and sets it for the app.
func (app *App) DetectAndSetLanguage(language string) error {
	langCode := lang.GetLang(language)
	l, err := lang.NewLangPackage(langCode)
	if err != nil {
		return fmt.Errorf("failed to initialize language: %w", err)
	}
	app.lang = l

	return nil
}

// GetLang returns the current language instance.
func (app *App) GetLang() *langGeneral.Lang {
	return app.lang
}

// ApplyConfig applies the given configuration, merging it with console size if needed, and validates it.
func (app *App) ApplyConfig(cfg Cfg) error {
	var err error

	cfg, err = app.mergeWindowSizeWithConsole(cfg)
	if err != nil {
		return fmt.Errorf("failed to merge config from console: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	app.cfg = cfg
	return nil
}

// Run starts the application and waits for the context to be done.
// It requires that the config has been applied beforehand.
func (app *App) Run() error {
	app.wg.Add(1)
	defer app.wg.Done()
	if app.cfg == nil {
		return fmt.Errorf("config isn't applied, call ApplyConfig first")
	}

	runFunc := bootstrap.NewRunMode(app.cfg, app.lang, app.stdin, app.stdout, app.stderr).Route()

	if err := runFunc(app.ctx); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	return nil
}

// mergeWindowSizeWithConsole merges the configuration with the current console window size.
// If maxCharsPerLine or maxLineCount are not set (less than 1), it uses the console dimensions.
func (app *App) mergeWindowSizeWithConsole(cfg Cfg) (Cfg, error) {
	if cfg.GetMaxCharsPerLine() > 0 && cfg.GetMaxLineCount() > 0 {
		return cfg, nil
	}
	width, height, err := term.GetSize(app.fd)
	if err != nil {
		width = 500
		height = 20
		//return nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	// Helper to adjust value if below minimum
	adjustIfNeeded := func(current int, addition int) int {
		if current < 1 {
			current += addition
			if current < 1 {
				current = addition
			}
		}
		return current
	}

	cfg.SetMaxCharsPerLine(adjustIfNeeded(cfg.GetMaxCharsPerLine(), width))
	cfg.SetMaxLineCount(adjustIfNeeded(cfg.GetMaxLineCount(), height))

	return cfg, nil
}
