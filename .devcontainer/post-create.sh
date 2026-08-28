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

arch="$(uname -m)"
case "${arch}" in
  x86_64) k8s_arch="amd64" ;;
  aarch64 | arm64) k8s_arch="arm64" ;;
  *) k8s_arch="${arch}" ;;
esac

echo "[devcontainer] Installing kubectl…"
kubectl_version="$(curl -fsSL https://dl.k8s.io/release/stable.txt)"
curl -fsSL -o /tmp/kubectl "https://dl.k8s.io/release/${kubectl_version}/bin/linux/${k8s_arch}/kubectl"
install -m 0755 /tmp/kubectl "${HOME}/.local/bin/kubectl"
rm -f /tmp/kubectl

echo "[devcontainer] Installing kind…"
kind_version="$(curl -fsSL https://api.github.com/repos/kubernetes-sigs/kind/releases/latest | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)"
kind_version="${kind_version:-v0.27.0}"
curl -fsSL -o "${HOME}/.local/bin/kind" "https://kind.sigs.k8s.io/dl/${kind_version}/kind-linux-${k8s_arch}"
chmod +x "${HOME}/.local/bin/kind"

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
echo "  kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client 2>/dev/null | head -n 1 || echo installed)"
echo "  kind:    $(kind version 2>/dev/null || echo installed)"
echo "  uv:      $(uv --version)"
echo "  specify: $(specify --version 2>/dev/null || specify version 2>/dev/null || echo installed)"
echo "  agent:   $(agent --version 2>/dev/null || echo 'installed — open a new shell if not found')"
echo ""
echo "[devcontainer] Cursor CLI: open a new shell, then run:  agent login"
echo "  (required to authenticate the agent CLI inside this Dev Container)"
echo "[devcontainer] Compose Postgres: DATABASE_HOST=host.docker.internal (same compose.yaml as the host path)"
