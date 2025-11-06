package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestOllamaClientGenerate(t *testing.T) {
	var captured struct {
		Path   string
		Method string
		Body   generateRequest
	}
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		captured.Path = req.URL.Path
		captured.Method = req.Method
		if err := json.NewDecoder(req.Body).Decode(&captured.Body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		body, _ := json.Marshal(generateResponse{Response: "answer", Done: true})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client := NewOllamaClient("http://unit-test", &http.Client{Transport: transport})
	resp, err := client.Generate(context.Background(), "model", "prompt")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if resp != "answer" {
		t.Fatalf("expected answer, got %q", resp)
	}
	if captured.Path != "/api/generate" || captured.Method != http.MethodPost {
		t.Fatalf("unexpected request: %#v", captured)
	}
	if captured.Body.Model != "model" || captured.Body.Prompt != "prompt" {
		t.Fatalf("unexpected body: %#v", captured.Body)
	}
}

func TestOllamaClientErrorStatus(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader("failure")),
			Header:     make(http.Header),
		}, nil
	})

	client := NewOllamaClient("http://unit-test", &http.Client{Transport: transport})
	_, err := client.Generate(context.Background(), "model", "prompt")
	if err == nil || !strings.Contains(err.Error(), "failure") {
		t.Fatalf("expected error containing failure, got %v", err)
	}
}
