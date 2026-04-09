# Unified multi-stage Dockerfile for Vite apps (admin / user)
# Usage:
#   docker build --build-arg APP_NAME=admin -f docker/web.Dockerfile .
#   docker build --build-arg APP_NAME=user  -f docker/web.Dockerfile .

# ---------------------------------------------------------------------------
# Stage 1 — Install dependencies (cached unless lockfile / package.json change)
# ---------------------------------------------------------------------------
FROM oven/bun:1.3.0 AS deps

WORKDIR /app/web

COPY web/package.json web/bun.lock web/turbo.json web/biome.json ./
COPY web/packages/ui/package.json ./packages/ui/
COPY web/packages/commitlint-config/package.json ./packages/commitlint-config/
COPY web/packages/typescript-config/package.json ./packages/typescript-config/
COPY web/apps/admin/package.json ./apps/admin/
COPY web/apps/user/package.json ./apps/user/

RUN bun install --frozen-lockfile

# ---------------------------------------------------------------------------
# Stage 2 — Build the target app via Turbo
# ---------------------------------------------------------------------------
FROM oven/bun:1.3.0 AS builder

ARG APP_NAME
RUN test -n "$APP_NAME" || (echo "ERROR: APP_NAME build arg is required (admin | user)" && exit 1)

WORKDIR /app/web

# Dependency tree from stage 1
COPY --from=deps /app/web/node_modules ./node_modules
COPY --from=deps /app/web/packages/ui/node_modules ./packages/ui/node_modules
COPY --from=deps /app/web/apps/${APP_NAME}/node_modules ./apps/${APP_NAME}/node_modules

# Source code — shared packages + target app only
COPY web/packages/ ./packages/
COPY web/apps/${APP_NAME}/ ./apps/${APP_NAME}/
COPY web/package.json web/turbo.json web/biome.json ./

RUN bun install --frozen-lockfile
RUN cd apps/${APP_NAME} && bun run build

# ---------------------------------------------------------------------------
# Stage 3 — Minimal production runtime (Vite static assets via Nginx)
# ---------------------------------------------------------------------------
FROM nginx:1.27-alpine AS runner

ARG APP_NAME
RUN test -n "$APP_NAME" || (echo "ERROR: APP_NAME build arg is required (admin | user)" && exit 1)

COPY --from=builder /app/web/apps/${APP_NAME}/dist /usr/share/nginx/html

RUN set -eu; \
  if [ "$APP_NAME" = "admin" ]; then \
    printf '%s\n' \
      'server {' \
      '  listen 3000;' \
      '  server_name _;' \
      '  root /usr/share/nginx/html;' \
      '  index index.html;' \
      '' \
      '  location = / {' \
      '    return 302 /admin/;' \
      '  }' \
      '' \
      '  location / {' \
      '    try_files $uri $uri/ $uri.html /admin/index.html;' \
      '  }' \
      '}' \
      > /etc/nginx/conf.d/default.conf; \
  elif [ "$APP_NAME" = "user" ]; then \
    printf '%s\n' \
      'server {' \
      '  listen 3000;' \
      '  server_name _;' \
      '  root /usr/share/nginx/html;' \
      '  index index.html;' \
      '' \
      '  location / {' \
      '    try_files $uri $uri/ $uri.html /index.html;' \
      '  }' \
      '}' \
      > /etc/nginx/conf.d/default.conf; \
  else \
    echo "ERROR: APP_NAME must be admin or user"; \
    exit 1; \
  fi

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
