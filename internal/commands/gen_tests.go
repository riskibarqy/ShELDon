package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewGenTestsCommand generates Go tests using an LLM backend.
func NewGenTestsCommand(deps Dependencies) *cobra.Command {
	var (
		file  string
		fn    string
		out   string
		model string
	)

	cmd := &cobra.Command{
		Use:   "gen-tests",
		Short: "Generate Go table-driven tests for a function",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(file) == "" || strings.TrimSpace(fn) == "" {
				return errors.New("--file and --func are required")
			}

			src, err := deps.Files.Read(file)
			if err != nil {
				return err
			}

			code := textutil.ExtractFunction(src, fn)
			if code == "" {
				return fmt.Errorf("function %s not found in %s", fn, file)
			}

			prompt := fmt.Sprintf("Write Go table-driven tests for this function. Use testing and testify. Keep names clear.\n\n%s", code)
			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			modelToUse := textutil.Choose(model, deps.Config.ModelReason)
			ans, err := deps.LLM.Generate(ctx, modelToUse, prompt)
			if err != nil {
				return err
			}

			if out == "" {
				out = fmt.Sprintf("%s_test.go", strings.ToLower(fn))
			}
			return deps.Files.WriteFile(out, textutil.NormalizeCode(ans))
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Go source file")
	cmd.Flags().StringVar(&fn, "func", "", "Function name")
	cmd.Flags().StringVar(&out, "out", "", "Output test filename (default <func>_test.go)")
	cmd.Flags().StringVar(&model, "model", "", "Override model (default OLDEV_MODEL_REASON)")
	return cmd
}
