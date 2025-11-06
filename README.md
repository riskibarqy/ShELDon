# ShELDon CLI

ShELDon is a single-binary CLI that wraps local LLM developer workflows powered by Ollama, structured with clean architecture and SOLID-oriented packages.

## Dependencies

- Go 1.21+ (tested on 1.25)
- An Ollama instance reachable from your machine
- Optional: locally downloaded models that match your `OLDEV_MODEL*` choices

## Installation

```bash
cp .env.example .env # optional overrides
go build ./cmd/oldev
```

## Usage

Display all commands (supports `--help`, `-h`, and `-help`):

```bash
./oldev --help
```

Override defaults for this invocation:

```bash
./oldev --model-general llama3.2:3b --timeout 90s gen-tests --file handlers/user.go --func CreateUser
```

Each sub-command still accepts its specific flags (e.g., `--model` on `gen-tests`, `--query` on `index-suggest`). Environment variables from `.env` fill in any values you omit.
Expect progress updates on stderr narrated by a particularly opinionated Sheldon Cooper—handy for tracking long-running requests (and for unsolicited life critiques).

## Testing

```bash
GOCACHE=$(mktemp -d) go test ./...
```

## Project Layout

- `cmd/oldev`: application entrypoint and wiring
- `internal/app`: Cobra root command and global configuration overrides
- `internal/commands`: use-case specific command handlers
- `internal/config`, `internal/llm`, `internal/system`, `internal/git`: infrastructure adapters
- `internal/textutil`, `internal/analysis`: shared utilities and domain helpers

Feel free to extend the CLI by adding new commands under `internal/commands` that lean on the existing abstractions for configuration, IO, and LLM access.
