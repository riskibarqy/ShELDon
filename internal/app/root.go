package app

import (
	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/commands"
)

// NewRootCommand wires all subcommands together and exposes global overrides.
func NewRootCommand(deps commands.Dependencies) *cobra.Command {
	cfg := deps.Config

	var (
		modelGeneral = cfg.ModelGeneral
		modelReason  = cfg.ModelReason
		modelCoder   = cfg.ModelCoder
		ollamaHost   = cfg.OllamaHost
		timeout      = cfg.Timeout
	)

	root := &cobra.Command{
		Use:   "oldev",
		Short: "Local-LLM dev tools powered by Ollama",
		Long:  "A Swiss Army CLI to supercharge Go backend work using local LLMs via Ollama.",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			flags := cmd.Root().PersistentFlags()
			if flags.Changed("model-general") {
				cfg.ModelGeneral = modelGeneral
			}
			if flags.Changed("model-reason") {
				cfg.ModelReason = modelReason
			}
			if flags.Changed("model-coder") {
				cfg.ModelCoder = modelCoder
			}
			if flags.Changed("ollama-host") {
				cfg.OllamaHost = ollamaHost
			}
			if flags.Changed("timeout") {
				cfg.Timeout = timeout
			}
			deps.Logger.Info(cmd, "Global parameters aligned. General=%s, Reason=%s, Coder=%s, Host=%s, Timeout=%s. Your welcome note may be sent later.",
				cfg.ModelGeneral, cfg.ModelReason, cfg.ModelCoder, cfg.OllamaHost, cfg.Timeout)
		},
	}

	root.PersistentFlags().StringVar(&modelGeneral, "model-general", cfg.ModelGeneral, "Default general-purpose LLM model")
	root.PersistentFlags().StringVar(&modelReason, "model-reason", cfg.ModelReason, "Default reasoning LLM model")
	root.PersistentFlags().StringVar(&modelCoder, "model-coder", cfg.ModelCoder, "Default coding LLM model")
	root.PersistentFlags().StringVar(&ollamaHost, "ollama-host", cfg.OllamaHost, "Ollama API host (e.g. http://localhost:11434)")
	root.PersistentFlags().DurationVar(&timeout, "timeout", cfg.Timeout, "LLM request timeout")

	root.AddCommand(
		commands.NewGenTestsCommand(deps),
		commands.NewCommitCommand(deps),
		commands.NewExplainAnalyzeCommand(deps),
		commands.NewPProfCommand(deps),
		commands.NewReviewMigrationCommand(deps),
		commands.NewCheckContractCommand(deps),
		commands.NewExplainLogsCommand(deps),
		commands.NewLintFixesCommand(deps),
		commands.NewGenK8sCommand(deps),
		commands.NewIndexSuggestCommand(deps),
		commands.NewPRReviewCommand(deps),
	)

	return root
}
