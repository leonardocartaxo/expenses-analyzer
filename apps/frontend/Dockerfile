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
EXPOSE 3000
CMD ["node", "apps/frontend/server.js"]
