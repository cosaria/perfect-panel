#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TOOLS_BIN="$ROOT_DIR/.tools/bin"

GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:-v2.11.4}"
GOIMPORTS_VERSION="${GOIMPORTS_VERSION:-v0.43.0}"

if ! command -v go >/dev/null 2>&1; then
  echo "错误：未检测到 Go，请先安装 Go 1.25+。"
  exit 1
fi

mkdir -p "$TOOLS_BIN"
export GOBIN="$TOOLS_BIN"

echo "安装 golangci-lint (${GOLANGCI_LINT_VERSION}) 到 $TOOLS_BIN"
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@"$GOLANGCI_LINT_VERSION"

echo "安装 goimports (${GOIMPORTS_VERSION}) 到 $TOOLS_BIN"
go install golang.org/x/tools/cmd/goimports@"$GOIMPORTS_VERSION"

test -x "$TOOLS_BIN/golangci-lint"
test -x "$TOOLS_BIN/goimports"

echo "Go 工具安装完成："
echo "  - $TOOLS_BIN/golangci-lint"
echo "  - $TOOLS_BIN/goimports"
echo "golangci-lint 版本：$("$TOOLS_BIN/golangci-lint" version | sed -n '1p')"
