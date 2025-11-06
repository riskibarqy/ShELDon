package commands

import (
	"context"
	"errors"

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
			if (path == "" || path == "-") && deps.Files.IsInteractive() {
				return errors.New("no lint findings provided; pipe golangci-lint output or pass --in <file>")
			}
			deps.Logger.Info(cmd, "Reading lint findings from %s. I do love enumerating flaws.", path)
			report, err := deps.Files.Read(path)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Captured %d bytes of contrition-worthy lint output.", len(report))

			prompt := "Given these golangci-lint findings, propose smallest code changes per issue. No broad refactors; targeted patches only.\n\n" + report
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelCoder)
			deps.Logger.Info(cmd, "Alerting model %s to prescribe minimal corrective surgery.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Remediation plan issued. Implement it before entropy wins.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to lint output or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
