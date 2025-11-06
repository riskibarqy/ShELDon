package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewReviewMigrationCommand reviews SQL migrations for safety.
func NewReviewMigrationCommand(deps Dependencies) *cobra.Command {
	var (
		path  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "review-migration",
		Short: "Review a Postgres migration for safety/downtime risks",
		RunE: func(cmd *cobra.Command, args []string) error {
			sql, err := deps.Files.Read(path)
			if err != nil {
				return err
			}

			prompt := "Review this Postgres migration for safety and downtime risk. Flag: full table rewrites, enum pitfalls, blocking DDL. Provide safer alternatives.\n\n" + sql
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

	cmd.Flags().StringVar(&path, "in", "-", "Path to SQL file or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
