SHELL := /bin/sh

.PHONY: bootstrap lint test dev server-bootstrap web-bootstrap server-lint web-lint server-test

bootstrap: server-bootstrap web-bootstrap

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