#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"
# shellcheck source=lib.sh
source "${ROOT}/scripts/local/lib.sh"

CLUSTER="${KIND_CLUSTER_NAME:-expenses-analyzer}"

ensure_kind_kubectl

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required for pnpm local:up (Docker Desktop on the Mac)" >&2
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
  if check_backend_health; then
    ready=1
    break
  fi
  sleep 2
done

if [[ "${ready}" -ne 1 ]]; then
  echo "backend /health never returned 200 { status: \"ok\" } (checked in-cluster)" >&2
  kubectl get pods >&2 || true
  exit 1
fi

# Printed for the host browser (Mac). If JetBrains/VS Code forwarded 8080/8081, stop those forwards.
echo "Frontend: http://127.0.0.1:8080"
echo "Backend:  http://127.0.0.1:8081"
