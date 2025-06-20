package main

import (
	"bytes"
	"context"
	"flag"
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

	// ========================
	// Add the domains to usage

	var usageBuf bytes.Buffer

	usageBuf.WriteString(fmt.Sprintf("\t%s: %s\n", "completion", "Generate completion script for the CLI"))
	for _, cfg := range configs {
		// Short usage for the domain
		domainUsage := fmt.Sprintf("\t%s: %s", cfg.Domain, cfg.Description)
		usageBuf.WriteString(domainUsage + "\n")

	}

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage of the CLI:")
		fmt.Fprintln(flag.CommandLine.Output(), "Available commands:")
		fmt.Fprintln(flag.CommandLine.Output(), usageBuf.String())
	}
	flag.Parse()

	// ========================

	switch os.Args[1] {
	case "completion":
		err = completion(ctx, configs)
		if err != nil {
			return fmt.Errorf("generating completion script: %w", err)
		}
	default:

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

	}

	return nil
}

// completion generates a completion script for the CLI.
func completion(ctx context.Context, configs map[string]Config) error {

	return nil
}
