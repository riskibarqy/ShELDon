package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewExplainLogsCommand diagnoses logs with the help of an LLM.
func NewExplainLogsCommand(deps Dependencies) *cobra.Command {
	var (
		path  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "explain-logs",
		Short: "Diagnose logs and propose next debugging steps",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps.Logger.Info(cmd, "Collecting logs from %s. Drama inevitably ensues.", path)
			logs, err := deps.Files.Read(path)
			if err != nil {
				return err
			}
			deps.Logger.Info(cmd, "Ingested %d bytes of operational angst.", len(logs))

			prompt := "You are an SRE. Diagnose cause and next steps from these logs. Return: Probable cause, Evidence lines, Next 3 commands to run.\n\n" + logs
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Model %s summoned to translate log-induced chaos into actionable steps.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Diagnosis dispatched. Please attempt not to break production again.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&path, "in", "-", "Path to logs or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
