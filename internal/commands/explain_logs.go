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
			logs, err := deps.Files.Read(path)
			if err != nil {
				return err
			}

			prompt := "You are an SRE. Diagnose cause and next steps from these logs. Return: Probable cause, Evidence lines, Next 3 commands to run.\n\n" + logs
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

	cmd.Flags().StringVar(&path, "in", "-", "Path to logs or '-' for stdin")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
