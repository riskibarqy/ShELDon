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

			query, err := deps.Files.Read(queryFile)
			if err != nil {
				return err
			}

			schema, err := deps.Shell.Run(schemaCmd)
			if err != nil && schemaCmd != "" {
				// Preserve original behaviour by ignoring failures but surfacing context.
				fmt.Fprintf(cmd.ErrOrStderr(), "schema command error: %v\n", err)
				schema = ""
			}

			prompt := "Suggest the ONE most impactful index for this query. Explain write amplification & size tradeoff.\n\nCurrent schema/indexes (optional):\n" + schema + "\n\nQuery:\n" + query
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

	cmd.Flags().StringVar(&schemaCmd, "schema-cmd", "", "Shell command to print schema/indexes (e.g. `psql -c \\d+ table`)")
	cmd.Flags().StringVar(&queryFile, "query", "", "Path to SQL file")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
