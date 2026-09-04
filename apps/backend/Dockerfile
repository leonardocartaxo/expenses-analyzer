FROM node:24-bookworm-slim
WORKDIR /repo
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc tsconfig.base.json ./
COPY apps/backend/package.json apps/backend/
COPY apps/frontend/package.json apps/frontend/
COPY packages/api-client/package.json packages/api-client/
RUN pnpm install --frozen-lockfile --filter @expenses/backend...
COPY apps/backend apps/backend
ENV PORT=3001
EXPOSE 3001
WORKDIR /repo/apps/backend
CMD ["pnpm", "start"]
