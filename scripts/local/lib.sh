#!/usr/bin/env bash
# Shared helpers for pnpm local:* (kind path).

# True when running inside a container (Dev Container / Docker-outside-of-Docker).
in_container() {
  [[ -f /.dockerenv ]] || [[ -f /run/.containerenv ]]
}

# Install kind + kubectl into ~/.local/bin when missing (Mac host or Dev Container).
# Uses the current OS/arch so the same script works on darwin and linux.
ensure_kind_kubectl() {
  export PATH="${HOME}/.local/bin:${PATH}"
  mkdir -p "${HOME}/.local/bin"

  local os arch
  case "$(uname -s)" in
    Darwin) os=darwin ;;
    Linux) os=linux ;;
    *)
      echo "unsupported OS for auto-install of kind/kubectl: $(uname -s)" >&2
      return 1
      ;;
  esac
  case "$(uname -m)" in
    x86_64 | amd64) arch=amd64 ;;
    aarch64 | arm64) arch=arm64 ;;
    *) arch="$(uname -m)" ;;
  esac

  if ! command -v kubectl >/dev/null 2>&1 || ! kubectl version --client >/dev/null 2>&1; then
    echo "Installing kubectl (${os}/${arch}) to ${HOME}/.local/bin…"
    local kubectl_version
    kubectl_version="$(curl -fsSL https://dl.k8s.io/release/stable.txt)"
    curl -fsSL -o "${HOME}/.local/bin/kubectl" \
      "https://dl.k8s.io/release/${kubectl_version}/bin/${os}/${arch}/kubectl"
    chmod +x "${HOME}/.local/bin/kubectl"
  fi

  if ! command -v kind >/dev/null 2>&1 || ! kind version >/dev/null 2>&1; then
    echo "Installing kind (${os}/${arch}) to ${HOME}/.local/bin…"
    local kind_version
    kind_version="$(curl -fsSL https://api.github.com/repos/kubernetes-sigs/kind/releases/latest | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)"
    kind_version="${kind_version:-v0.27.0}"
    curl -fsSL -o "${HOME}/.local/bin/kind" \
      "https://kind.sigs.k8s.io/dl/${kind_version}/kind-${os}-${arch}"
    chmod +x "${HOME}/.local/bin/kind"
  fi
}

# kind writes kubeconfig with https://127.0.0.1:<port>. That works on the host;
# from a Dev Container, 127.0.0.1 is the workspace container — use host.docker.internal.
# API server cert may not include that name on older clusters, so skip TLS verify only in DooD.
fix_kind_kubeconfig_for_dood() {
  local cluster="${1:?cluster name}"
  local context="kind-${cluster}"
  local server fixed

  if ! in_container; then
    return 0
  fi

  if ! kubectl config get-contexts "${context}" >/dev/null 2>&1; then
    kind export kubeconfig --name "${cluster}" >/dev/null
  fi

  server="$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"${context}\")].cluster.server}")"
  if [[ -z "${server}" ]]; then
    return 0
  fi

  fixed="${server//127.0.0.1/host.docker.internal}"
  fixed="${fixed//localhost/host.docker.internal}"
  # Always re-apply when in DooD: server may already be host.docker.internal from a prior run
  # without insecure-skip-tls-verify (API cert typically lacks that SAN until cluster recreate).
  if [[ "${fixed}" == *host.docker.internal* ]]; then
    kubectl config set-cluster "${context}" --server="${fixed}" --insecure-skip-tls-verify=true >/dev/null
  fi
}

# Health check via the backend pod (avoids host:8081 conflicts with IDE port forwards
# and Docker Desktop hairpin from the Dev Container).
check_backend_health() {
  kubectl exec deploy/backend -- node -e '
fetch("http://127.0.0.1:3001/health")
  .then(async (res) => {
    const body = await res.json().catch(() => ({}));
    process.exit(res.status === 200 && body.status === "ok" ? 0 : 1);
  })
  .catch(() => process.exit(1));
' >/dev/null 2>&1
}

# Optional: also try published URL when LOCAL_HEALTH_URL is set (host debugging).
check_health_url() {
  local url="${1:?}"
  node -e '
fetch(process.argv[1])
  .then(async (res) => {
    const body = await res.json().catch(() => ({}));
    process.exit(res.status === 200 && body.status === "ok" ? 0 : 1);
  })
  .catch(() => process.exit(1));
' "${url}"
}
