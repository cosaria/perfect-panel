# Unified multi-stage Dockerfile for Next.js apps (admin / user)
# Usage:
#   docker build --build-arg APP_NAME=admin -f docker/web.Dockerfile .
#   docker build --build-arg APP_NAME=user  -f docker/web.Dockerfile .

# ---------------------------------------------------------------------------
# Stage 1 — Install dependencies (cached unless lockfile / package.json change)
# ---------------------------------------------------------------------------
FROM oven/bun:1.3.0 AS deps

WORKDIR /app

COPY package.json bun.lock turbo.json ./
COPY packages/ui/package.json ./packages/ui/
COPY packages/commitlint-config/package.json ./packages/commitlint-config/
COPY packages/typescript-config/package.json ./packages/typescript-config/
COPY apps/admin/package.json ./apps/admin/
COPY apps/user/package.json ./apps/user/

RUN bun install --frozen-lockfile

# ---------------------------------------------------------------------------
# Stage 2 — Build the target app via Turbo
# ---------------------------------------------------------------------------
FROM oven/bun:1.3.0 AS builder

ARG APP_NAME
RUN test -n "$APP_NAME" || (echo "ERROR: APP_NAME build arg is required (admin | user)" && exit 1)

WORKDIR /app

# Dependency tree from stage 1
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages/ui/node_modules ./packages/ui/node_modules
COPY --from=deps /app/apps/${APP_NAME}/node_modules ./apps/${APP_NAME}/node_modules

# Source code — shared packages + target app only
COPY packages/ ./packages/
COPY apps/${APP_NAME}/ ./apps/${APP_NAME}/
COPY package.json turbo.json biome.json ./

ENV NEXT_TELEMETRY_DISABLED=1

RUN bunx turbo run build --filter=ppanel-${APP_NAME}-web

# ---------------------------------------------------------------------------
# Stage 3 — Minimal production runtime (Node.js Alpine)
# ---------------------------------------------------------------------------
FROM node:22-alpine AS runner

ARG APP_NAME

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Next.js standalone output + static assets + public directory
COPY --from=builder /app/apps/${APP_NAME}/.next/standalone ./
COPY --from=builder /app/apps/${APP_NAME}/.next/static ./apps/${APP_NAME}/.next/static
COPY --from=builder /app/apps/${APP_NAME}/public ./apps/${APP_NAME}/public

RUN chown -R nextjs:nodejs /app

USER nextjs

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production
ENV APP_NAME=${APP_NAME}

EXPOSE 3000

# Shell form so $APP_NAME is expanded at runtime
CMD node apps/${APP_NAME}/server.js
