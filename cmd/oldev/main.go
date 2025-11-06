package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/riskiramdan/ShELDon/internal/app"
	"github.com/riskiramdan/ShELDon/internal/commands"
	"github.com/riskiramdan/ShELDon/internal/config"
	"github.com/riskiramdan/ShELDon/internal/git"
	"github.com/riskiramdan/ShELDon/internal/llm"
	"github.com/riskiramdan/ShELDon/internal/system"
)

func main() {
	// Allow running as `go run . --help` without Cobra parsing flag.CommandLine twice.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	_ = flag.CommandLine.Parse([]string{})

	cfg := config.Load(config.OSEnvReader{})

	deps := commands.Dependencies{
		Config: cfg,
		LLM:    llm.NewOllamaClient(cfg.OllamaHost, http.DefaultClient),
		Files:  system.NewOSFileManager(os.Stdin),
		Git:    git.CLIClient{},
		Shell:  system.BashShell{},
	}

	root := app.NewRootCommand(deps)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
