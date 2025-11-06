package llm

import "context"

// Client describes the behaviour required from any large-language-model backend.
type Client interface {
	Generate(ctx context.Context, model, prompt string) (string, error)
}
