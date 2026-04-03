FROM golang:1.25-alpine AS builder

ARG TARGETARCH
ARG VERSION
ENV CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH}

RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN BUILD_TIME=$(date -u +"%Y-%m-%d %H:%M:%S") && \
    go build -ldflags="-s -w -X 'github.com/perfect-panel/server/pkg/constant.Version=${VERSION}' -X 'github.com/perfect-panel/server/pkg/constant.BuildTime=${BUILD_TIME}'" -o /app/ppanel ppanel.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/ppanel /app/ppanel

EXPOSE 8080

ENTRYPOINT ["/app/ppanel"]
CMD ["run", "--config", "/app/etc/ppanel.yaml"]
