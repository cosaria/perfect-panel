SHELL := /bin/sh

.PHONY: bootstrap lint test dev build format typecheck clean \
	server-bootstrap web-bootstrap server-lint web-lint server-test \
	server-build web-build server-format web-format web-typecheck \
	server-clean web-clean embed embed-admin embed-user build-all

bootstrap: server-bootstrap web-bootstrap
	@command -v lefthook >/dev/null 2>&1 && lefthook install || echo "Warning: lefthook not found. Install via: brew install lefthook"

server-bootstrap:
	cd server && go mod download

web-bootstrap:
	CI=true bun install

lint: server-lint web-lint

server-lint:
	cd server && golangci-lint run && go vet ./...

web-lint:
	bun run lint

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
		(bun run dev:$(APP)) & \
		wait

build: server-build web-build

server-build:
	cd server && go build -o bin/ppanel-server .

web-build:
	bun run build

format: server-format web-format

server-format:
	cd server && go fmt ./... && goimports -w .

web-format:
	bun run format

typecheck: web-typecheck

web-typecheck:
	bun run typecheck

embed-admin:
	bun run build --filter=ppanel-admin-web
	rm -rf server/web/admin-dist/*
	cp -r apps/admin/out/* server/web/admin-dist/

embed-user:
	bun run build --filter=ppanel-user-web
	rm -rf server/web/user-dist/*
	cp -r apps/user/out/* server/web/user-dist/

embed: embed-admin embed-user

build-all: embed
	cd server && go build -tags embed -ldflags="-s -w" -o bin/ppanel .

clean: server-clean web-clean

server-clean:
	rm -rf server/bin/ server/web/admin-dist/* server/web/user-dist/*
	touch server/web/admin-dist/.gitkeep server/web/user-dist/.gitkeep

web-clean:
	rm -rf apps/*/.next apps/*/.turbo .turbo
