# All-in-one build: admin Vite SPA + user Vite SPA embedded in Go binary
# Usage: docker build -t ppanel .

# Stage 1: Build both frontends
FROM oven/bun:1-alpine AS web-builder
WORKDIR /app
COPY web/package.json web/bun.lock web/turbo.json web/biome.json ./web/
COPY web/packages/ ./web/packages/
COPY web/apps/admin/ ./web/apps/admin/
COPY web/apps/user/ ./web/apps/user/
WORKDIR /app/web
RUN bun install --frozen-lockfile
RUN bun run build --filter=ppanel-admin-web --filter=ppanel-user-web

# Stage 2: Build Go binary with embedded frontends
FROM golang:1.25-alpine AS builder
ARG TARGETARCH
ARG VERSION=dev
ENV CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH}

WORKDIR /build
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
COPY --from=web-builder /app/web/apps/admin/dist ./web/admin-dist
COPY --from=web-builder /app/web/apps/user/dist ./web/user-dist

RUN BUILD_TIME=$(date -u +"%Y-%m-%d %H:%M:%S") && \
    go build -tags embed \
      -ldflags="-s -w -X 'github.com/perfect-panel/server/config.Version=${VERSION}' -X 'github.com/perfect-panel/server/config.BuildTime=${BUILD_TIME}'" \
      -o /app/ppanel ppanel.go

# Stage 3: Minimal runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ppanel /app/ppanel
COPY --from=builder /build/cache ./cache

EXPOSE 8080
ENTRYPOINT ["/app/ppanel"]
CMD ["run", "--config", "/app/etc/ppanel.yaml"]
