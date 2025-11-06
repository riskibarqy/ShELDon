# ShELDon CLI

ShELDon is a single-binary CLI that wraps local LLM developer workflows powered by Ollama.

## Getting Started

```bash
cp .env.example .env # optional overrides
go build ./cmd/oldev
./oldev --help
```

Set the `OLLAMA_HOST` and `OLDEV_MODEL*` variables to point at your preferred models. Default values are documented in `.env.example`.

## Testing

```bash
GOCACHE=$(mktemp -d) go test ./...
```

The repository is structured following clean architecture and SOLID principles with dedicated packages for configuration, infrastructure (LLM, git, shell, filesystem), and the command use-cases.
