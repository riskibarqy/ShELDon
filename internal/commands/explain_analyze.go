package commands

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewExplainAnalyzeCommand explains PostgreSQL EXPLAIN ANALYZE output.
func NewExplainAnalyzeCommand(deps Dependencies) *cobra.Command {
	var (
		path  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "explain-analyze",
		Short: "Explain a PostgreSQL EXPLAIN ANALYZE plan and suggest indexes/rewrite",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (path == "" || path == "-") && deps.Files.IsInteractive() {
				return errors.New("no plan provided; pipe EXPLAIN ANALYZE output or use --in <file>")
			}
			deps.Logger.Info(cmd, "Acquiring EXPLAIN ANALYZE output from %s. I hope it brought a bibliography.", path)
			plan, err := deps.Files.Read(path)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Digesting a modest %d bytes of planner musings.", len(plan))

			prompt := "Explain the PostgreSQL EXPLAIN ANALYZE below. Give: 1) bottlenecks, 2) missing/misused indexes, 3) rewrite suggestion.\n\n" + plan
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Deploying model %s to interpret the planner's cryptic opera.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Diagnosis rendered. If databases could blush, this one just did.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to plan file or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
