package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {

	if len(os.Args) < 2 {
		return fmt.Errorf("no command provided")
	}

	// Read configs from the directory
	configs := make(map[string]Config)
	err := filepath.WalkDir("dir", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		cfg, err := readConfig(path)
		if err != nil {
			return fmt.Errorf("reading config from %s: %w", path, err)
		}

		configs[cfg.Domain] = cfg

		return nil
	})
	if err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	fmt.Printf("%+v\n", configs)

	// ========================

	cfg, ok := configs[os.Args[1]]
	if !ok {
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}

	args := os.Args[2:]

	command, err := cfg.Fs(args)
	if err != nil {
		return fmt.Errorf("getting command from config: %w", err)
	}

	// Execute the command

	shell := "/bin/bash"
	cmd := exec.CommandContext(ctx, shell, "-c", command.Cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
