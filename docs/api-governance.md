# OpenAPI Governance

This document defines the contract-governance policy for the documented server APIs exported from [`server/cmd/openapi.go`](/Users/admin/Codes/ProxyCode/perfect-panel/server/cmd/openapi.go).

## Governed Surfaces

- `docs/openapi/admin.json`
- `docs/openapi/common.json`
- `docs/openapi/user.json`

These artifacts are generated from Huma-declared first-party JSON APIs and are the only surfaces governed by Phase 5B.

## Naming Policy

- `operationId` stays unique, URL-safe, and stable once consumed by generated clients.
- Tags must be declared at the top level with a description.
- Summaries should use short imperative or descriptive sentence case without leaking internal implementation jargon.
- Shared authenticated JSON APIs should reuse the Phase 5 RFC 9457 problem-details error contract.

## Security And Error Responses

- Documented authenticated APIs must declare explicit bearer-token security requirements.
- Documented first-party JSON APIs must expose explicit common `4xx` responses in OpenAPI.
- Generated specs must remain consumable by `@hey-api/openapi-ts` for both `web/apps/admin` and `web/apps/user`.

## Excluded Surfaces

- webhook and callback protocols such as payment notifications and Telegram webhook handlers
- node polling and conditional-cache endpoints
- redirect-based OAuth callback flows
- init/bootstrap setup endpoints in `server/initialize/`

These surfaces are real HTTP endpoints, but they are excluded from this OpenAPI governance layer because they do not behave like ordinary documented JSON APIs.

## Breaking Change Workflow

1. Update route declarations and shared conventions in `server/routers/`.
2. Regenerate specs and generated clients with `bun run openapi`.
3. Review generated diff in `docs/openapi/`, `web/apps/admin/services/`, and `web/apps/user/services/`.
4. Call out intentional contract renames or compatibility risks in review before merge.
