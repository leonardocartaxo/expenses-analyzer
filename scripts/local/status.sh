#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=lib.sh
source "${ROOT}/scripts/local/lib.sh"

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"
HEALTH_URL="${LOCAL_HEALTH_URL:-$(default_health_url)}"

if ! command -v kind >/dev/null 2>&1 || ! command -v kubectl >/dev/null 2>&1; then
  echo "kind and kubectl are required for pnpm local:status" >&2
  exit 1
fi

if ! kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
  echo "kind cluster ${CLUSTER}: not running"
  exit 1
fi

fix_kind_kubeconfig_for_dood "${CLUSTER}"

echo "kind cluster ${CLUSTER}: running"
kubectl get pods -o wide
if node -e "fetch(process.argv[1]).then(async (res) => { const body = await res.json().catch(() => ({})); console.log('health', res.status, JSON.stringify(body)); process.exit(res.status === 200 && body.status === 'ok' ? 0 : 1); }).catch((err) => { console.error(err.message); process.exit(1); })" "${HEALTH_URL}"; then
  echo "stack: ready"
else
  echo "stack: not ready"
  exit 1
fi
