#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=lib.sh
source "${ROOT}/scripts/local/lib.sh"

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"

ensure_kind_kubectl

if kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
  kind delete cluster --name "${CLUSTER}"
fi
