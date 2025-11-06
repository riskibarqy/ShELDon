package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewCommitCommand generates a Conventional Commit message using LLM support.
func NewCommitCommand(deps Dependencies) *cobra.Command {
	var model string

	cmd := &cobra.Command{
		Use:   "llm-commit",
		Short: "Generate a Conventional Commit message from staged changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			diff, err := deps.Git.Diff("--staged")
			if err != nil {
				return err
			}
			if strings.TrimSpace(diff) == "" {
				return errors.New("no staged changes")
			}

			prompt := "Write a concise Conventional Commit message for this diff. One-line summary, then bullets of key changes.\n\n" + diff
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			ans, err := deps.LLM.Generate(ctx, textutil.Choose(model, deps.Config.ModelGeneral), prompt)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(ans))
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Override model (default OLDEV_MODEL)")
	return cmd
}
