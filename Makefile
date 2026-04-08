SHELL := /bin/sh

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
TOOLS_BIN := $(ROOT_DIR)/.tools/bin
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOIMPORTS := $(TOOLS_BIN)/goimports
INSTALL_GO_TOOLS := $(ROOT_DIR)/.github/scripts/install-go-tools.sh
BUN_INSTALL_FLAGS ?=

.PHONY: bootstrap tools lint test dev build format typecheck clean \
	ensure-go-tools server-bootstrap web-bootstrap server-lint web-lint server-test \
	server-build web-build server-format web-format web-typecheck \
	server-clean web-clean embed embed-admin embed-user build-all

bootstrap: tools server-bootstrap web-bootstrap
	@command -v lefthook >/dev/null 2>&1 && lefthook install || echo "Warning: lefthook not found. Install via: brew install lefthook"

tools:
	sh "$(INSTALL_GO_TOOLS)"

ensure-go-tools:
	@test -x "$(GOLANGCI_LINT)" || { echo "Missing $(GOLANGCI_LINT). Run 'make tools' from the repo root."; exit 1; }
	@test -x "$(GOIMPORTS)" || { echo "Missing $(GOIMPORTS). Run 'make tools' from the repo root."; exit 1; }

server-bootstrap:
	cd server && go mod download

web-bootstrap:
	cd web && CI=true bun install $(BUN_INSTALL_FLAGS)

lint: server-lint web-lint

server-lint: ensure-go-tools
	cd server && "$(GOLANGCI_LINT)" run && go vet ./...

web-lint:
	cd web && bun run lint

test: server-test

server-test:
	cd server && go test ./...

dev:
	@if [ -z "$(APP)" ]; then echo "Usage: make dev APP=admin|user"; exit 1; fi
	@case "$(APP)" in \
		admin|user) ;; \
		*) echo "APP must be admin or user"; exit 1 ;; \
	esac
	@trap 'kill 0' INT TERM EXIT; \
		(cd server && go run . run --config etc/ppanel.yaml) & \
		(cd web && bun run dev:$(APP)) & \
		wait

build: server-build web-build

server-build:
	cd server && go build -o bin/ppanel-server .

web-build:
	cd web && bun run build

format: server-format web-format

server-format: ensure-go-tools
	cd server && go fmt ./... && "$(GOIMPORTS)" -w .

web-format:
	cd web && bun run format

typecheck: web-typecheck

web-typecheck:
	cd web && bun run typecheck

embed-admin:
	cd web && bun run build --filter=ppanel-admin-web
	rm -rf server/web/admin-dist/*
	cp -r web/apps/admin/dist/* server/web/admin-dist/

embed-user:
	cd web && bun run build --filter=ppanel-user-web
	rm -rf server/web/user-dist/*
	cp -r web/apps/user/out/* server/web/user-dist/

embed: embed-admin embed-user

build-all: embed
	cd server && go build -tags embed -ldflags="-s -w" -o bin/ppanel .

clean: server-clean web-clean

server-clean:
	rm -rf server/bin/ server/web/admin-dist/* server/web/user-dist/*
	touch server/web/admin-dist/.gitkeep server/web/user-dist/.gitkeep

web-clean:
	rm -rf web/apps/*/.next web/apps/*/.turbo web/.turbo
