#!/usr/bin/env bash
set -euo pipefail

export PATH="${HOME}/.local/bin:${PATH}"

append_path() {
  local file="$1"
  if [[ -f "${file}" ]] || touch "${file}"; then
    if ! grep -q '\.local/bin' "${file}" 2>/dev/null; then
      printf '\n# local tools (uv, cursor agent, …)\nexport PATH="%s/.local/bin:$PATH"\n' "${HOME}" >> "${file}"
    fi
  fi
}

append_path "${HOME}/.zshrc"
append_path "${HOME}/.bashrc"
append_path "${HOME}/.profile"

echo "[devcontainer] Installing uv…"
curl -LsSf https://astral.sh/uv/install.sh | sh

echo "[devcontainer] Installing specify-cli (Spec Kit)…"
uv tool install specify-cli

echo "[devcontainer] Installing Cursor CLI (agent)…"
curl https://cursor.com/install -fsS | bash

echo "[devcontainer] Ready."
echo "  node:    $(node --version)"
echo "  npm:     $(npm --version)"
echo "  uv:      $(uv --version)"
echo "  specify: $(specify --version 2>/dev/null || specify version 2>/dev/null || echo installed)"
echo "  agent:   $(agent --version 2>/dev/null || echo 'installed — open a new shell if not found')"
echo ""
echo "[devcontainer] Cursor CLI: open a new shell, then run:  agent login"
echo "  (required to authenticate the agent CLI inside this Dev Container)"
