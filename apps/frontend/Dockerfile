FROM node:24-bookworm-slim AS builder
WORKDIR /repo
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc tsconfig.base.json ./
COPY packages/api-client packages/api-client
COPY apps/frontend apps/frontend
RUN pnpm install --frozen-lockfile --filter @expenses/frontend...
ENV NEXT_TELEMETRY_DISABLED=1
RUN pnpm --filter @expenses/frontend build

FROM node:24-bookworm-slim
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
ENV API_BASE_URL=http://backend:3001
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=builder /repo/apps/frontend/.next/standalone ./
COPY --from=builder /repo/apps/frontend/.next/static ./apps/frontend/.next/static
# Next 16.3 + Node 24 NFT may copy only @swc/helpers/esm under next's nested
# node_modules (missing package.json exports). Replace with the full package.
RUN set -eux; \
  helpers="$(find node_modules/.pnpm -path '*/node_modules/@swc/helpers/package.json' | head -n1 | xargs -r dirname)"; \
  next_helpers="$(find node_modules/.pnpm -path '*/next@*/node_modules/@swc/helpers' -type d | head -n1)"; \
  if [ -n "${helpers}" ] && [ -n "${next_helpers}" ]; then \
    rm -rf "${next_helpers}"; \
    cp -a "${helpers}" "${next_helpers}"; \
  fi
EXPOSE 3000
CMD ["node", "apps/frontend/server.js"]
