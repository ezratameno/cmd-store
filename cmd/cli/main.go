package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
)

var version = "dev"
var commit = "unknown"
var date = "unknown"

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

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home directory: %w", err)
	}

	configsDir := filepath.Join(home, ".cmd-store")

	if _, err := os.Stat(configsDir); os.IsNotExist(err) {
		return fmt.Errorf("configs directory does not exist: %s", configsDir)
	}

	// Read configs from the directory
	configs := make(map[string]Config)
	err = filepath.WalkDir(configsDir, func(path string, d os.DirEntry, err error) error {
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
	case "version":
		fmt.Println("CLI Version:", version)
		fmt.Println("Commit:", commit)
		fmt.Println("Date:", date)
	case "completion":
		err = completion(configsDir, configs)
		if err != nil {
			return fmt.Errorf("generating completion script: %w", err)
		}
	default:

		cfg, ok := configs[os.Args[1]]
		if !ok {
			return fmt.Errorf("unknown command: %s", os.Args[1])
		}

		args := os.Args[2:]

		command, opts, err := cfg.Fs(args)
		if err != nil {
			return fmt.Errorf("getting command from config: %w", err)
		}

		if opts.DryMode {
			// If dry mode is enabled, print the command and return
			fmt.Println("Dry run mode enabled. Command will not be executed.")
			fmt.Printf("Command: %s\n", command.Cmd)
			return nil
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

//go:embed templates
var templateFile embed.FS

// completion generates a completion script for the CLI.
func completion(configsDir string, configs map[string]Config) error {
	programPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	programPath = filepath.Base(programPath)

	tmplFile := "templates/completion.sh.tmpl"

	type Domain struct {
		Name string
		Cmds string
	}
	type TemplateData struct {
		ProgramName string
		Domains     []Domain
		DomainsStr  string
	}

	var data TemplateData
	data.ProgramName = programPath

	for _, cfg := range configs {
		var domain Domain
		domain.Name = cfg.Domain
		for _, cmd := range cfg.Commands {
			domain.Cmds = fmt.Sprintf("%s %s", domain.Cmds, cmd.Name)
		}
		data.Domains = append(data.Domains, domain)

		data.DomainsStr += fmt.Sprintf("%s ", cfg.Domain)
	}

	tmpl, err := template.ParseFS(templateFile, tmplFile)
	if err != nil {
		return fmt.Errorf("parsing template file %s: %w", tmplFile, err)
	}

	completionLoc := filepath.Join(configsDir, "cmd-store-completion.sh")
	file, err := os.Create(completionLoc)
	if err != nil {
		return fmt.Errorf("creating completion file %s: %w", completionLoc, err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	fmt.Println("Completion script generated at:", completionLoc)
	fmt.Println("To enable completion for the current shell, run:")
	fmt.Printf("\tsource %s\n", completionLoc)

	fmt.Println("You can also add the following line to your shell configuration file (e.g., ~/.bashrc or ~/.zshrc):")
	fmt.Printf("\tif [ -f %s ]; then\n", completionLoc)
	fmt.Printf("\t\tsource %s\n", completionLoc)
	fmt.Println("\tfi")

	// var buf bytes.Buffer
	// var domains []string

	// buf.WriteString("#! /bin/bash\n")
	// for _, cfg := range configs {
	// 	domains = append(domains, cfg.Domain)

	// 	var cmds []string
	// 	for _, cmd := range cfg.Commands {
	// 		cmds = append(cmds, cmd.Name)
	// 	}

	// 	// Add the cmds to the completion script
	// 	buf.WriteString(fmt.Sprintf("complete -W \"%s\" %s %s\n", strings.Join(cmds, " "), programPath, cfg.Domain))
	// }

	// // Add the domain commands to the completion script
	// buf.WriteString(fmt.Sprintf("complete -W \"%s\" %s\n", strings.Join(domains, " "), programPath))

	// // for _, cfg := range configs {

	// // }
	// fmt.Println(buf.String())
	return nil
}
