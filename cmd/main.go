package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	app "tail/internal/app"
	config "tail/internal/config"
	internalContext "tail/internal/context"
	"tail/internal/formatter"
)

var (
	defaultStderr = os.Stderr
	defaultStdout = os.Stdout
	defaultStdin  = os.Stdin
)

func main() {
	err := execute()
	if err != nil {
		_, _ = fmt.Fprintf(defaultStderr, "Error: %v\n", err)
	}
	os.Exit(exitCodeFromError(err))
}

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	var exitErr formatter.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}

	return 1
}

func execute() (retErr error) {
	var (
		err         error
		cfg         app.Cfg
		application *app.App
		wg          sync.WaitGroup
		ctx, cancel = internalContext.NewMainContext(context.Background(), defaultStderr, &wg)
		appCancel   func() error
	)

	defer func() {
		cancel()
		if appCancel != nil {
			if errC := appCancel(); errC != nil {
				retErr = errors.Join(retErr, errC)
			}
		}
		wg.Wait()
	}()

	// Create application
	application, appCancel, err = app.NewApp(os.Stdout.Fd(), ctx, &wg)
	if err != nil {
		return err
	}

	// Set default standard streams
	application.SetDefaultStd(defaultStdin, defaultStdout, defaultStderr)

	// Detect language package
	if err = application.DetectAndSetLanguage(os.Getenv("LANG")); err != nil {
		return err
	}

	// Get Config
	cfg, err = config.ReadConf(application.GetLang())
	if err != nil {
		return err
	}

	// Apply config
	err = application.ApplyConfig(cfg)
	if err != nil {
		return err
	}

	// Run application
	err = application.Run()
	if err != nil {
		return err
	}

	return nil

}
