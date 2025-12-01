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
			cmd.SilenceUsage = true

			deps.Logger.Info(cmd, "Evaluating your staged diff. Contain your anticipation.")
			diff, err := deps.Git.Diff("--staged")
			if err != nil {
				return err
			}
			if strings.TrimSpace(diff) == "" {
				return errors.New("no staged changes")
			}

			// Defensive: cap diff size and escape triple backticks to avoid block echoes.
			const maxDiffLen = 10000
			if len(diff) > maxDiffLen {
				diff = diff[len(diff)-maxDiffLen:]
				deps.Logger.Info(cmd, "Diff truncated to last %d bytes for model input.", maxDiffLen)
			}
			diff = strings.ReplaceAll(diff, "```", "`​``") // insert zero-width char to break triple backticks

			basePrompt := fmt.Sprintf(`Write ONLY a single-line Conventional Commit message for the diff below.
Format exactly as "<type(scope)?: >concise summary in lowercase present tense".
Keep the line at or below %d characters—be concise instead of adding follow-up text.
Do not include bullets, explanations, reviews, or multiple lines. Return just the commit header without quotes.`, deps.Config.MaxSummaryLen)

			prompt := basePrompt + "\n\n" + diff
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelUse := textutil.Choose(model, deps.Config.ModelGeneral)
			deps.Logger.Info(cmd, "Summoning model %s to translate chaos into convention.", modelUse)

			// Try up to N times: generate -> normalize -> validate -> if fails, retry with stricter prompt.
			const maxAttempts = 3
			var lastCandidate string
			for attempt := 1; attempt <= maxAttempts; attempt++ {
				ans, err := deps.LLM.Generate(ctx, modelUse, prompt)
				if err != nil {
					return err
				}
				candidate := applyPrefix(normalizeCommitMessage(ans), prefix)
				deps.Logger.Info(cmd, "LLM returned candidate: %s", candidate)

				// Keep only the first line (we want single-line summary for commit header).
				firstLine := strings.SplitN(candidate, "\n", 2)[0]
				firstLine = strings.TrimSpace(firstLine)

				// If too long, optionally ask the model to shorten (demonstrated by helper).
				if len(firstLine) > deps.Config.MaxSummaryLen {
					deps.Logger.Info(cmd, "Candidate summary too long (%d chars), requesting shortening.", len(firstLine))
					short, err := shortenSummaryWithLLM(ctx, deps, modelUse, firstLine)
					if err == nil && short != "" {
						firstLine = short
					}
				}

				if validateConventionalCommit(firstLine) {
					// success
					fmt.Fprintln(cmd.OutOrStdout(), firstLine)
					deps.Logger.Info(cmd, "Commit message prepared. Praise can be mailed to apartment 4A.")
					if autoCommit {
						deps.Logger.Info(cmd, "Executing git commit with the freshly minted prose.")
						if err := deps.Git.Commit(firstLine); err != nil {
							return err
						}
						deps.Logger.Info(cmd, "Commit recorded. I recommend celebratory string theory.")
					}
					return nil
				}

				deps.Logger.Info(cmd, "Candidate did not match conventional-collected rules (attempt %d).", attempt)
				lastCandidate = firstLine
				// make prompt stricter for next attempt
				prompt = basePrompt + "\n\n" +
					"The previous candidate was invalid. Produce a single-line Conventional Commit summary only. " +
					"Use one of the types: feat, fix, docs, style, refactor, perf, test, chore. " +
					"Example: feat(parser): handle edge case\n\n" + diff
			}

			// all attempts failed -> surface last candidate for manual editing
			deps.Logger.Info(cmd, "All LLM attempts failed to produce a valid Conventional Commit message.")
			fmt.Fprintln(cmd.OutOrStdout(), lastCandidate)
			return errors.New("failed to generate a valid conventional commit message; please edit manually")
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Override model (default SHELDON_MODEL)")
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
		bannedPrefixes := []string{
			"this is a code review",
			"overall,",
			"overall:",
		}
		skip := false
		for _, prefix := range bannedPrefixes {
			if strings.HasPrefix(lower, prefix) {
				skip = true
				break
			}
		}
		if skip {
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

// validateConventionalCommit checks a simple Conventional Commit pattern.
// Adjust allowed types / scope characters to suit your repo policy.
func validateConventionalCommit(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	candidate := s
	for candidate != "" {
		if CompiledConventionalCommit.MatchString(candidate) {
			return true
		}
		sep := strings.IndexAny(candidate, " \t")
		if sep == -1 {
			break
		}
		candidate = strings.TrimSpace(candidate[sep+1:])
	}
	return false
}

// shortenSummaryWithLLM asks the model to shorten a one-line summary.
// This is a best-effort helper that returns shortened string or error.
func shortenSummaryWithLLM(ctx context.Context, deps Dependencies, model, long string) (string, error) {
	prompt := `Shorten the following Conventional Commit summary to <=` + fmt.Sprintf("%d", deps.Config.MaxSummaryLen) + ` characters without changing meaning.
Return only the shortened single-line summary.` + "\n\n" + long
	ans, err := deps.LLM.Generate(ctx, model, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.SplitN(ans, "\n", 2)[0]), nil
}
