#!/usr/bin/env bash
# Shared helpers for pnpm local:* (kind path).

# True when running inside a container (Dev Container / Docker-outside-of-Docker).
in_container() {
  [[ -f /.dockerenv ]] || [[ -f /run/.containerenv ]]
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

# Health URL used by scripts (wait/status). Host browser URLs stay 127.0.0.1.
default_health_url() {
  if in_container; then
    echo "http://host.docker.internal:8081/health"
  else
    echo "http://127.0.0.1:8081/health"
  fi
}
