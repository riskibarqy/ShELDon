package commands

import (
	"context"
	"errors"

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
			if (path == "" || path == "-") && deps.Files.IsInteractive() {
				return errors.New("no migration provided; supply --in <file> or pipe SQL text")
			}
			deps.Logger.Info(cmd, "Opening SQL migration %s. May the DDL be ever in your favor.", path)
			sql, err := deps.Files.Read(path)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Catalogued %d characters of schema meddling.", len(sql))

			prompt := "Review this Postgres migration for safety and downtime risk. Flag: full table rewrites, enum pitfalls, blocking DDL. Provide safer alternatives.\n\n" + sql
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Consulting model %s for a pre-flight safety inspection.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Migration risk report delivered. Proceed, cautiously, if at all.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to SQL file or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
