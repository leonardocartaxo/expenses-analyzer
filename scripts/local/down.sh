#!/usr/bin/env bash
set -euo pipefail

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required for pnpm local:down" >&2
  exit 1
fi

if kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
  kind delete cluster --name "${CLUSTER}"
fi
