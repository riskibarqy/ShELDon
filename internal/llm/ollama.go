package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OllamaClient implements Client by calling the Ollama HTTP API.
type OllamaClient struct {
	host       string
	httpClient *http.Client
}

// NewOllamaClient builds an Ollama-backed LLM client.
func NewOllamaClient(host string, httpClient *http.Client) *OllamaClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &OllamaClient{
		host:       strings.TrimRight(host, "/"),
		httpClient: httpClient,
	}
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate sends a prompt to Ollama's /api/generate endpoint.
func (c *OllamaClient) Generate(ctx context.Context, model, prompt string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("marshal generate request: %w", err)
	}

	url := c.host + "/api/generate"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		buf := new(bytes.Buffer)
		_, readErr := buf.ReadFrom(resp.Body)
		if readErr != nil {
			return "", fmt.Errorf("ollama error status %d", resp.StatusCode)
		}
		return "", fmt.Errorf("ollama error: %s", strings.TrimSpace(buf.String()))
	}

	var out generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode generate response: %w", err)
	}
	return out.Response, nil
}
