package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewPProfCommand analyses pprof output and provides guidance.
func NewPProfCommand(deps Dependencies) *cobra.Command {
	var (
		path  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "pprof-analyze",
		Short: "Analyze a pprof -top output and suggest concrete optimizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			text, err := deps.Files.Read(path)
			if err != nil {
				return err
			}

			prompt := "You're a Go perf engineer. Analyze this pprof -top and say EXACTLY which funcs to attack and how (allocs, pools, JSON, etc.).\n\n" + text
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

	cmd.Flags().StringVar(&path, "in", "-", "Path to pprof text or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
