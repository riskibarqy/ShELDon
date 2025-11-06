package commands

import (
	"context"

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
			plan, err := deps.Files.Read(path)
			if err != nil {
				return err
			}

			prompt := "Explain the PostgreSQL EXPLAIN ANALYZE below. Give: 1) bottlenecks, 2) missing/misused indexes, 3) rewrite suggestion.\n\n" + plan
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			ans, err := deps.LLM.Generate(ctx, textutil.Choose(model, deps.Config.ModelGeneral), prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to plan file or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
