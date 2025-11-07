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
	var (
		model      string
		prefix     string
		autoCommit bool
	)

	cmd := &cobra.Command{
		Use:   "llm-commit",
		Short: "Generate a Conventional Commit message from staged changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps.Logger.Info(cmd, "Evaluating your staged diff. Contain your anticipation.")
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

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Summoning model %s to translate chaos into convention.", modelUse)
			ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
			if err != nil {
				return err
			}

			message := applyPrefix(normalizeCommitMessage(ans), prefix)
			fmt.Fprintln(cmd.OutOrStdout(), message)
			deps.Logger.Info(cmd, "Commit message prepared. Praise can be mailed to apartment 4A.")

			if autoCommit {
				deps.Logger.Info(cmd, "Executing git commit with the freshly minted prose.")
				if err := deps.Git.Commit(message); err != nil {
					return err
				}
				deps.Logger.Info(cmd, "Commit recorded. I recommend celebratory string theory.")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Override model (default OLDEV_MODEL)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Text prepended to the first line of the generated commit message")
	cmd.Flags().BoolVar(&autoCommit, "autocommit", false, "If true, automatically run git commit with the generated message")
	return cmd
}

func applyPrefix(message, prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || message == "" {
		return message
	}
	parts := strings.SplitN(message, "\n", 2)
	parts[0] = fmt.Sprintf("%s %s", prefix, strings.TrimSpace(parts[0]))
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts, "\n")
}

func normalizeCommitMessage(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	lines := strings.Split(raw, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(cleaned) == 0 || cleaned[len(cleaned)-1] == "" {
				continue
			}
			cleaned = append(cleaned, "")
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "here is") || strings.HasPrefix(lower, "here's") {
			continue
		}
		if strings.HasPrefix(trimmed, "```") || strings.HasSuffix(trimmed, "```") {
			continue
		}
		if strings.HasPrefix(trimmed, "**") && strings.HasSuffix(trimmed, "**") && len(trimmed) > 4 {
			trimmed = strings.Trim(trimmed, "*")
			trimmed = strings.TrimSpace(trimmed)
		}
		cleaned = append(cleaned, trimmed)
	}
	for len(cleaned) > 0 && cleaned[len(cleaned)-1] == "" {
		cleaned = cleaned[:len(cleaned)-1]
	}
	return strings.Join(cleaned, "\n")
}
