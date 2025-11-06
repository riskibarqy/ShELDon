package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewLintFixesCommand proposes minimal patches based on golangci-lint output.
func NewLintFixesCommand(deps Dependencies) *cobra.Command {
	var (
		path  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "lint-fixes",
		Short: "Propose smallest code changes for golangci-lint findings",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := deps.Files.Read(path)
			if err != nil {
				return err
			}

			prompt := "Given these golangci-lint findings, propose smallest code changes per issue. No broad refactors; targeted patches only.\n\n" + report
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			ans, err := deps.LLM.Generate(ctx, textutil.Choose(model, deps.Config.ModelCoder), prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to lint output or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
