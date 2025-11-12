package config

import (
	"os"
	"time"
)

// Config groups all runtime settings required by the CLI.
type Config struct {
	ModelGeneral string
	ModelReason  string
	ModelCoder   string
	OllamaHost   string
	Timeout      time.Duration
}

// EnvReader abstracts environment variable access to support testing.
type EnvReader interface {
	LookupEnv(string) (string, bool)
}

// OSEnvReader is the production implementation backed by os.LookupEnv.
type OSEnvReader struct{}

// LookupEnv delegates to os.LookupEnv.
func (OSEnvReader) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

const (
	defaultModelGeneral = "llama3.1:8b"
	defaultModelReason  = "deepseek-r1:7b"
	defaultModelCoder   = "qwen2.5-coder:1.5b"
	defaultOllamaHost   = "http://localhost:11434"
	defaultTimeout      = 120 * time.Second
	envModelGeneral     = "SHELDON_MODEL"
	envModelReason      = "SHELDON_MODEL_REASON"
	envModelCoder       = "SHELDON_MODEL_CODER"
	envOllamaHost       = "OLLAMA_HOST"
	envTimeout          = "SHELDON_TIMEOUT"
)

// Load builds a Config using environment variables with sensible defaults.
func Load(reader EnvReader) Config {
	cfg := Config{
		ModelGeneral: valueOrDefault(reader, envModelGeneral, defaultModelGeneral),
		ModelReason:  valueOrDefault(reader, envModelReason, defaultModelReason),
		ModelCoder:   valueOrDefault(reader, envModelCoder, defaultModelCoder),
		OllamaHost:   valueOrDefault(reader, envOllamaHost, defaultOllamaHost),
		Timeout:      defaultTimeout,
	}

	if str, ok := reader.LookupEnv(envTimeout); ok && str != "" {
		if duration, err := time.ParseDuration(str); err == nil {
			cfg.Timeout = duration
		}
	}

	return cfg
}

func valueOrDefault(reader EnvReader, key, def string) string {
	if value, ok := reader.LookupEnv(key); ok && value != "" {
		return value
	}
	return def
}
