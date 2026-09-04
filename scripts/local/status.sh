#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=lib.sh
source "${ROOT}/scripts/local/lib.sh"

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"

ensure_kind_kubectl

if ! kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
  echo "kind cluster ${CLUSTER}: not running"
  exit 1
fi

fix_kind_kubeconfig_for_dood "${CLUSTER}"

echo "kind cluster ${CLUSTER}: running"
kubectl get pods -o wide
if check_backend_health; then
  echo "health 200 {\"status\":\"ok\"} (in-cluster)"
  echo "stack: ready"
else
  echo "stack: not ready"
  exit 1
fi
