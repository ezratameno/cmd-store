package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Commands    []Command `yaml:"commands"`
	Description string    `yaml:"description"`

	// The name of the command, for example 'aws'.
	Domain string
}

type Command struct {
	Cmd         string `yaml:"cmd"`
	Description string `yaml:"description"`

	// The short name of the command used for execution of the command
	Name string `yaml:"name"`
}

func readConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshalling config file: %w", err)
	}

	domain := filepath.Base(path)
	domain = strings.TrimSuffix(domain, filepath.Ext(domain))

	config.Domain = domain

	return config, nil
}

// Fs creates a flag set for the given config and parses the arguments.
// And return the command to execute
func (c Config) Fs(args []string) (Command, error) {
	domainFs := flag.NewFlagSet(c.Domain, flag.ExitOnError)

	// Register the available commands in the usage function
	var domainUsageBuf bytes.Buffer
	for _, cmd := range c.Commands {

		// ==========================
		// Domain

		// Short usage for the domain
		cmdDomainUsage := fmt.Sprintf("\t%s: %s", cmd.Name, cmd.Description)
		domainUsageBuf.WriteString(cmdDomainUsage + "\n")

	}

	domainFs.Usage = func() {
		fmt.Fprintln(domainFs.Output(), "Usage of", c.Domain+":")
		fmt.Fprintln(domainFs.Output(), "Available commands:")
		fmt.Fprintln(domainFs.Output(), domainUsageBuf.String())
	}

	domainFs.Parse(args)

	// Check for sub commands

	if len(args) < 1 {
		return Command{}, fmt.Errorf("no command provided")
	}

	for _, cmd := range c.Commands {

		if cmd.Name != args[0] {
			continue
		}

		// =============================
		// Command flag
		cmdFs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

		cmdUsage := fmt.Sprintf("\tDescription: %s \n \tCommand: %s", cmd.Description, cmd.Cmd)

		cmdFs.Usage = func() {
			fmt.Fprintln(cmdFs.Output(), "Usage of", cmd.Name+":")
			fmt.Fprintln(cmdFs.Output(), cmdUsage)
			cmdFs.PrintDefaults()
		}

		cmdFs.Parse(args[1:])

		return cmd, nil

	}

	return Command{}, fmt.Errorf("unknown command: %s", args[0])

}
