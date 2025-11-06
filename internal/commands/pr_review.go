package commands

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewPRReviewCommand runs an LLM-powered review for the current branch diff.
func NewPRReviewCommand(deps Dependencies) *cobra.Command {
	var (
		base  string
		model string
	)

	cmd := &cobra.Command{
		Use:   "pr-review",
		Short: "Structured code review for the current branch diff",
		RunE: func(cmd *cobra.Command, args []string) error {
			if base == "" {
				base = "origin/main"
			}

			deps.Logger.Info(cmd, "Calculating diff against %s. Time to expose questionable decisions.", base)
			diff, err := deps.Git.Diff(base + "...HEAD")
			if err != nil {
				return err
			}
			if strings.TrimSpace(diff) == "" {
				return errors.New("no diff vs base")
			}

			prompt := "Code review with 5 sections: Correctness, Complexity, Style, Tests, Security. Be specific, cite file:line. Keep under 200 lines.\n\n" + diff
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelReason)
			deps.Logger.Info(cmd, "Deploying model %s to perform a code review that actually reads the diff.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(ans))
			if err == nil {
				deps.Logger.Info(cmd, "Review complete. Remember, sarcasm is my love language.")
			}
			return err
		},
	}

	cmd.Flags().StringVar(&base, "base", "origin/main", "Base ref for diff")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
