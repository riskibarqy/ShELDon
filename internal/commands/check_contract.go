package commands

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/analysis"
	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewCheckContractCommand finds mismatches between API spec and implementation.
func NewCheckContractCommand(deps Dependencies) *cobra.Command {
	var (
		specPath string
		implDir  string
		model    string
	)

	cmd := &cobra.Command{
		Use:   "check-contract",
		Short: "Find mismatches between API spec (proto/OpenAPI) and Go handlers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if specPath == "" {
				return errors.New("--spec is required")
			}

			spec, err := deps.Files.Read(specPath)
			if err != nil {
				return err
			}

			scanner := analysis.NewImplScanner(deps.Files)
			snippets, err := scanner.Scan(implDir)
			if err != nil {
				return err
			}

			prompt := "Find mismatches between this API spec and Go handlers. Report missing fields, wrong types, status codes, pagination rules.\n\nSPEC:\n" + spec + "\n\nIMPL SNIPPETS:\n" + snippets
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			ans, err := deps.LLM.Generate(ctx, textutil.Choose(model, deps.Config.ModelReason), prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			return err
		},
	}

	cmd.Flags().StringVar(&specPath, "spec", "", "Path to API spec file (proto/openapi)")
	cmd.Flags().StringVar(&implDir, "impl", ".", "Directory with Go handlers")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	_ = cmd.MarkFlagRequired("spec")
	return cmd
}
