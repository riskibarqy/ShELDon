package app

import (
	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/commands"
)

// NewRootCommand wires all subcommands together.
func NewRootCommand(deps commands.Dependencies) *cobra.Command {
	root := &cobra.Command{
		Use:   "oldev",
		Short: "Local-LLM dev tools powered by Ollama",
		Long:  "A Swiss Army CLI to supercharge Go backend work using local LLMs via Ollama.",
	}

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
