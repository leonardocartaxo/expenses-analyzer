#!/usr/bin/env bash
set -euo pipefail

export PATH="${HOME}/.local/bin:${PATH}"
mkdir -p "${HOME}/.local/bin"

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

echo "[devcontainer] Enabling Corepack / pnpm…"
sudo corepack enable
corepack prepare pnpm@10.15.1 --activate

# Docker-outside-of-Docker: CLI is in the image; socket comes from the host mount.
if [[ -S /var/run/docker.sock ]]; then
  sock_gid="$(stat -c '%g' /var/run/docker.sock)"
  if ! getent group "${sock_gid}" >/dev/null 2>&1; then
    sudo groupadd --gid "${sock_gid}" docker-host
  fi
  sock_group="$(getent group "${sock_gid}" | cut -d: -f1)"
  sudo usermod -aG "${sock_group}" node || true
  # Same-session access without re-login (local Dev Container only).
  if ! docker info >/dev/null 2>&1; then
    sudo chmod 666 /var/run/docker.sock || true
  fi
  echo "[devcontainer] docker.sock mounted (gid ${sock_gid} / group ${sock_group})"
else
  echo "[devcontainer] WARNING: /var/run/docker.sock is missing."
  echo "  Start Docker Desktop on the host, then rebuild this Dev Container."
fi

echo "[devcontainer] Installing kind + kubectl (same helper as pnpm local:*)…"
# shellcheck source=../scripts/local/lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/scripts/local/lib.sh"
ensure_kind_kubectl

echo "[devcontainer] Installing uv…"
curl -LsSf https://astral.sh/uv/install.sh | sh

echo "[devcontainer] Installing specify-cli (Spec Kit)…"
uv tool install specify-cli

echo "[devcontainer] Installing Cursor CLI (agent)…"
curl https://cursor.com/install -fsS | bash

echo "[devcontainer] Ready."
echo "  node:    $(node --version)"
echo "  npm:     $(npm --version)"
echo "  pnpm:    $(pnpm --version)"
echo "  docker:  $(docker --version 2>/dev/null || echo 'missing — rebuild image')"
echo "  compose: $(docker compose version 2>/dev/null || echo 'missing — rebuild image')"
echo "  kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null | head -n 1 || echo installed)"
echo "  kind:    $(kind version 2>/dev/null || echo installed)"
echo "  uv:      $(uv --version)"
echo "  specify: $(specify --version 2>/dev/null || specify version 2>/dev/null || echo installed)"
echo "  agent:   $(agent --version 2>/dev/null || echo 'installed — open a new shell if not found')"
echo ""
echo "[devcontainer] Cursor CLI: open a new shell, then run:  agent login"
echo "  (required to authenticate the agent CLI inside this Dev Container)"
echo "[devcontainer] Compose Postgres: DATABASE_HOST=${DATABASE_HOST:-unset} (expect host.docker.internal; same compose.yaml as host path)"
echo "[devcontainer] After rebuild: docker compose up -d  (uses host Docker Desktop via docker.sock)"
if [[ "${DATABASE_HOST:-}" != "host.docker.internal" ]]; then
  echo "[devcontainer] WARNING: DATABASE_HOST is not host.docker.internal."
  echo "  export DATABASE_HOST=host.docker.internal   # then restart Nest"
fi
