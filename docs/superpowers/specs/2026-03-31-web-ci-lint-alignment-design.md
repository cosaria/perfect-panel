# Web CI Lint Alignment Design

## Goal

Make GitHub Actions validate the web workspace using the same Biome-only lint entrypoint used locally, so shared packages such as `@workspace/ui` are covered by CI.

## Current Problem

- The current workflow runs `bun run lint` only inside `web/apps/admin` and `web/apps/user`.
- Shared workspace packages under `web/packages/*` are not included in those app-local lint commands.
- The real workspace lint entrypoint is `web/package.json -> turbo lint`, which already validates `@workspace/ui`, `ppanel-admin-web`, and `ppanel-user-web` together.

## Chosen Approach

Add a dedicated `web-lint` GitHub Actions job that runs `bun run lint` from the `web` root.

Keep `web-admin-validate` and `web-user-validate` for app-specific `check-types` and `build` verification, but remove their per-app lint steps.

## Why This Approach

- Matches the current repository tooling contract after the Biome migration.
- Catches shared-package formatting and lint issues in CI.
- Minimizes workflow churn by preserving existing app build and typecheck jobs.
- Keeps CI output granular enough to distinguish lint failures from app build/typecheck failures.

## Workflow Changes

- Add `web-lint` job gated by the existing `changes` outputs.
- Reuse the local `setup-bun` action for dependency installation.
- Run `bun run lint` with `working-directory: web`.
- Remove the app-local lint steps from admin and user jobs.
- Add the new job to `monorepo-summary` needs and summary output.

## Validation

- Run `actionlint` against the workflow.
- Run `./script/check-workflow-contract.sh`.
