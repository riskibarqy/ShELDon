package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewIndexSuggestCommand proposes the most impactful index for a query.
func NewIndexSuggestCommand(deps Dependencies) *cobra.Command {
	var (
		schemaCmd string
		queryFile string
		model     string
	)

	cmd := &cobra.Command{
		Use:   "index-suggest",
		Short: "Suggest the single most impactful index for a Postgres query",
		RunE: func(cmd *cobra.Command, args []string) error {
			if queryFile == "" {
				return errors.New("--query is required")
			}

			deps.Logger.Info(cmd, "Loading query from %s. I trust it follows first normal form.", queryFile)
			query, err := deps.Files.Read(queryFile)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Query length: %d characters. Suitable for academic peer review.", len(query))

			deps.Logger.Info(cmd, "If a schema command exists, I shall execute it with geologic punctuality.")
			schema, err := deps.Shell.Run(schemaCmd)
			if err != nil && schemaCmd != "" {
				// Preserve original behaviour by ignoring failures but surfacing context.
				fmt.Fprintf(cmd.ErrOrStderr(), "schema command error: %v\n", err)
				schema = ""
			} else if schemaCmd != "" {
				deps.Logger.Info(cmd, "Schema details acquired. I now know more about your database than HR does about you.")
			}

			prompt := "Suggest the ONE most impactful index for this query. Explain write amplification & size tradeoff.\n\nCurrent schema/indexes (optional):\n" + schema + "\n\nQuery:\n" + query
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Asking model %s to identify the mathematically optimal index.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Index advice delivered. Apply it before the optimizer files a complaint.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&schemaCmd, "schema-cmd", "", "Shell command to print schema/indexes (e.g. `psql -c \\d+ table`)")
	cmd.Flags().StringVar(&queryFile, "query", "", "Path to SQL file")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
