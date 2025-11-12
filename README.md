# ShELDon CLI

ShELDon is a single-binary CLI that wraps local LLM developer workflows powered by Ollama, structured with clean architecture and SOLID-oriented packages.

## Dependencies

- Go 1.21+ (tested on 1.25)
- An Ollama instance reachable from your machine
- Optional: locally downloaded models that match your `SHELDON_MODEL*` choices

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

## Command Examples

- **`gen-tests`** – generate table-driven tests for a Go function  
  ```bash
  sheldon gen-tests --file internal/service/user.go --func CreateUser --out create_user_test.go
  ```

- **`llm-commit`** – produce a Conventional Commit message from staged changes  
  ```bash
  git add .
  sheldon llm-commit --prefix "feat(api):" --autocommit
  ```
  `--prefix` prepends text to the first line, while `--autocommit` tells the CLI to immediately run `git commit -m`.

- **`explain-analyze`** – interpret a PostgreSQL execution plan  
  ```bash
  psql -d dbname -c "EXPLAIN (ANALYZE, BUFFERS) SELECT ..." | sheldon explain-analyze --in -
  ```

- **`pprof-analyze`** – review `pprof -top` output for optimizations  
  ```bash
  go tool pprof -top cpu.pprof > /tmp/pprof.txt
  sheldon pprof-analyze --in /tmp/pprof.txt
  ```

- **`review-migration`** – assess migration safety  
  ```bash
  sheldon review-migration --in migrations/20241001_add_column.sql
  ```

- **`check-contract`** – detect spec vs. handler mismatches  
  ```bash
  sheldon check-contract --spec api/openapi.yaml --impl internal/handlers
  ```

- **`explain-logs`** – diagnose log output and suggest next actions  
  ```bash
  tail -n 500 logs/app.log | sheldon explain-logs --in -
  ```

- **`lint-fixes`** – summarize minimal fixes for golangci-lint findings  
  ```bash
  golangci-lint run ./... --out-format tab | sheldon lint-fixes --in -
  ```

- **`gen-k8s`** – scaffold Kubernetes Deployment + HPA YAML  
  ```bash
  sheldon gen-k8s --app users-api --port 8080 --min 2 --max 6 --out deploy.yaml
  ```

- **`index-suggest`** – recommend a single impactful index  
  ```bash
  sheldon index-suggest --query queries/slow.sql --schema-cmd "psql -d dbname -c '\d+ users'"
  ```

- **`pr-review`** – run an LLM code review against a diff  
  ```bash
  sheldon pr-review --base origin/main
  ```

- **`completion`** – generate shell completions (bash|zsh|fish|powershell)  
  ```bash
  sheldon completion zsh > "${fpath[1]}/_sheldon"
  ```

## Testing

```bash
GOCACHE=$(mktemp -d) go test ./...
```

## Project Layout

- `cmd/sheldon`: application entrypoint and wiring
- `internal/app`: Cobra root command and global configuration overrides
- `internal/commands`: use-case specific command handlers
- `internal/config`, `internal/llm`, `internal/system`, `internal/git`: infrastructure adapters
- `internal/textutil`, `internal/analysis`: shared utilities and domain helpers

Feel free to extend the CLI by adding new commands under `internal/commands` that lean on the existing abstractions for configuration, IO, and LLM access.
