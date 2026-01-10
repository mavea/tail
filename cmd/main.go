package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"tail/internal/adapters/renderers"
	"tail/internal/adapters/sources"
	"tail/internal/adapters/targets"
	"tail/internal/config"
	"tail/internal/processor"
)

func main() {
	// Get Config
	cfg, err := config.ReadConf()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Make Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init source
	in, err := sources.NewScanner(ctx, cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Init target
	target, err := targets.New(cfg, bufio.NewWriter(os.Stdout))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Init renderer
	out, cancelOut, err := renderers.New(ctx, target, cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer cancelOut()

	// Init processor
	service := processor.NewService(in, out)

	// Start processor
	if err = service.Run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
