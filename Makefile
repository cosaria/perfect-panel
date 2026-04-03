SHELL := /bin/sh

.PHONY: bootstrap lint test dev build format typecheck clean \
	server-bootstrap web-bootstrap server-lint web-lint server-test \
	server-build web-build server-format web-format web-typecheck \
	server-clean web-clean

bootstrap: server-bootstrap web-bootstrap
	@command -v lefthook >/dev/null 2>&1 && lefthook install || echo "Warning: lefthook not found. Install via: brew install lefthook"

server-bootstrap:
	cd server && go mod download

web-bootstrap:
	cd web && CI=true bun install

lint: server-lint web-lint

server-lint:
	cd server && golangci-lint run && go vet ./...

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
		(cd web && bun run dev --filter=ppanel-$(APP)-web) & \
		wait

build: server-build web-build

server-build:
	cd server && go build -o bin/ppanel-server .

web-build:
	cd web && bun run build

format: server-format web-format

server-format:
	cd server && go fmt ./... && goimports -w .

web-format:
	cd web && bun run format

typecheck: web-typecheck

web-typecheck:
	cd web && bun run typecheck

clean: server-clean web-clean

server-clean:
	rm -rf server/bin/

web-clean:
	cd web && rm -rf apps/*/.next apps/*/.turbo .turbo