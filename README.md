# ShELDon CLI

ShELDon is a single-binary CLI that wraps local LLM developer workflows powered by Ollama, structured with clean architecture and SOLID-oriented packages.

## Dependencies

- Go 1.21+ (tested on 1.25)
- An Ollama instance reachable from your machine
- Optional: locally downloaded models that match your `OLDEV_MODEL*` choices

## Installation

```bash
cp .env.example .env # optional overrides
make build           # builds ./bin/sheldon
```

## Usage

Display all commands (supports `--help`, `-h`, and `-help`):

```bash
./bin/sheldon --help
```

Override defaults for this invocation:

```bash
./bin/sheldon --model-general llama3.2:3b --timeout 90s gen-tests --file handlers/user.go --func CreateUser
```

Each sub-command still accepts its specific flags (e.g., `--model` on `gen-tests`, `--query` on `index-suggest`). Environment variables from `.env` fill in any values you omit.
Expect progress updates on stderr narrated by a particularly opinionated Sheldon Cooper—handy for tracking long-running requests (and for unsolicited life critiques).

### Global Installation

- `make install` installs the binary to your `GOBIN`/`GOPATH/bin` as `sheldon`, enabling `sheldon --help` from any directory.
- `make shell-alias` appends an alias to `~/.zshrc` and `~/.bashrc` pointing at the local build (`./bin/sheldon`), if you prefer not to install globally.
- Reload your shell (e.g., `source ~/.zshrc`) after running either command to activate the new executable/alias.

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
