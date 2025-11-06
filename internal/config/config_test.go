package config

import (
	"testing"
	"time"
)

type fakeEnv map[string]string

func (f fakeEnv) LookupEnv(key string) (string, bool) {
	v, ok := f[key]
	return v, ok
}

func TestLoadDefaults(t *testing.T) {
	cfg := Load(fakeEnv{})
	if cfg.ModelGeneral != defaultModelGeneral {
		t.Fatalf("expected default ModelGeneral %q, got %q", defaultModelGeneral, cfg.ModelGeneral)
	}
	if cfg.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout %s, got %s", defaultTimeout, cfg.Timeout)
	}
}

func TestLoadOverrides(t *testing.T) {
	env := fakeEnv{
		envModelGeneral: "custom-general",
		envModelReason:  "reason",
		envModelCoder:   "coder",
		envOllamaHost:   "http://example",
		envTimeout:      "30s",
	}

	cfg := Load(env)
	if cfg.ModelGeneral != "custom-general" {
		t.Fatalf("expected ModelGeneral override, got %q", cfg.ModelGeneral)
	}
	if cfg.Timeout != 30*time.Second {
		t.Fatalf("expected timeout override 30s, got %s", cfg.Timeout)
	}
}
