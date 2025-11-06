package commands

import (
	"context"
	"errors"

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
			if (path == "" || path == "-") && deps.Files.IsInteractive() {
				return errors.New("no pprof data supplied; pipe output or pass --in <file>")
			}
			deps.Logger.Info(cmd, "Preparing to interpret pprof output from %s. Performance sins, reveal yourselves.", path)
			text, err := deps.Files.Read(path)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Parsed %d bytes of flame fodder. Science commences.", len(text))

			prompt := "You're a Go perf engineer. Analyze this pprof -top and say EXACTLY which funcs to attack and how (allocs, pools, JSON, etc.).\n\n" + text
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelReason)
			deps.Logger.Info(cmd, "Engaging model %s for a performance autopsy.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Optimization guidance broadcast. Your CPU just sent a thank-you card.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to pprof text or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
