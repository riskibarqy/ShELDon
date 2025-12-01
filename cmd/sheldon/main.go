package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/riskiramdan/ShELDon/internal/app"
	"github.com/riskiramdan/ShELDon/internal/commands"
	"github.com/riskiramdan/ShELDon/internal/config"
	"github.com/riskiramdan/ShELDon/internal/git"
	"github.com/riskiramdan/ShELDon/internal/llm"
	"github.com/riskiramdan/ShELDon/internal/logging"
	"github.com/riskiramdan/ShELDon/internal/system"
)

func main() {
	normalizeHelpFlag()

	// Allow running as `go run . --help` without Cobra parsing flag.CommandLine twice.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	_ = flag.CommandLine.Parse([]string{})

	// Load overrides from .env when present; ignore missing file errors.
	_ = godotenv.Load(".env")

	cfg := config.Load(config.OSEnvReader{})

	deps := commands.Dependencies{
		Config: &cfg,
		LLM:    llm.NewOllamaClient(cfg.OllamaHost, http.DefaultClient),
		Files:  system.NewOSFileManager(os.Stdin),
		Git:    git.CLIClient{},
		Shell:  system.BashShell{},
		Logger: logging.NewSheldonLogger(),
	}

	root := app.NewRootCommand(deps)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// normalizeHelpFlag maps the single-dash help alias to Cobra's expected --help.
func normalizeHelpFlag() {
	for i := range os.Args {
		if os.Args[i] == "-help" {
			os.Args[i] = "--help"
		}
	}
}
