#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"
# shellcheck source=lib.sh
source "${ROOT}/scripts/local/lib.sh"

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"
HEALTH_URL="${LOCAL_HEALTH_URL:-$(default_health_url)}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required for pnpm local:up" >&2
  exit 1
fi
if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is required for pnpm local:up" >&2
  exit 1
fi
if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required for pnpm local:up" >&2
  exit 1
fi

if ! kind get clusters 2>/dev/null | grep -qx "${CLUSTER}"; then
  kind create cluster --name "${CLUSTER}" --config deploy/kind/cluster.yaml
fi

fix_kind_kubeconfig_for_dood "${CLUSTER}"

docker build -t expenses-analyzer-backend:local -f apps/backend/Dockerfile .
docker build -t expenses-analyzer-frontend:local -f apps/frontend/Dockerfile .
kind load docker-image expenses-analyzer-backend:local --name "${CLUSTER}"
kind load docker-image expenses-analyzer-frontend:local --name "${CLUSTER}"
kubectl apply -k deploy/kustomize/overlays/kind

ready=0
for _ in $(seq 1 90); do
  if node -e "fetch(process.argv[1]).then(async (res) => { const body = await res.json().catch(() => ({})); process.exit(res.status === 200 && body.status === 'ok' ? 0 : 1); }).catch(() => process.exit(1))" "${HEALTH_URL}"; then
    ready=1
    break
  fi
  sleep 2
done

if [[ "${ready}" -ne 1 ]]; then
  echo "GET ${HEALTH_URL} never returned 200 { status: \"ok\" }" >&2
  kubectl get pods >&2 || true
  exit 1
fi

# Printed for the host browser (Mac), even when scripts run in the Dev Container.
echo "Frontend: http://127.0.0.1:8080"
echo "Backend:  http://127.0.0.1:8081"
