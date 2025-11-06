#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 1 ]; then
  echo "usage: scripts/add-shell-alias.sh /path/to/sheldon" >&2
  exit 1
fi

BIN_PATH="$1"
if [ ! -x "$BIN_PATH" ]; then
  echo "error: binary not executable: $BIN_PATH" >&2
  exit 1
fi

BIN_PATH="$(cd "$(dirname "$BIN_PATH")" && pwd)/$(basename "$BIN_PATH")"
ALIAS_CMD="alias sheldon=\"${BIN_PATH}\""
SHELL_FILES=("$HOME/.zshrc" "$HOME/.bashrc")

for rc in "${SHELL_FILES[@]}"; do
  if [ -f "$rc" ]; then
    if grep -Fq "$ALIAS_CMD" "$rc"; then
      echo "Sheldon alias already present in $rc"
    else
      {
        echo ""
        echo "# Added by ShELDon setup on $(date)"
        echo "$ALIAS_CMD"
      } >>"$rc"
      echo "Appended Sheldon alias to $rc"
    fi
  else
    echo "Skipping $rc (file not found)"
  fi
done

echo "Reload your shell or run 'source ~/.zshrc' (or ~/.bashrc) to activate the alias."
