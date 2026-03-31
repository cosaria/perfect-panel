# Root GitHub Control Plane

This repository now treats `.github/` as the single control plane for CI, release automation, issue hygiene, and contributor entrypoints.

## Routing Overview

```text
pull_request
    |
    +--> monorepo-check
            |
            +--> workflow-validate
            +--> server-validate when server/ or shared root files change
            +--> web-admin-validate when web/ or shared root files change
            +--> web-user-validate when web/ or shared root files change
            +--> monorepo-summary (the only required check)

push: main
    |
    +--> monorepo-check full sweep

push: main|next|beta
    |
    +--> web-release
            |
            +--> semantic-release emits web-v* tags only

release: published
    |
    +--> web-publish-release-assets when release tag starts with web-v

push: v*
    |
    +--> server-release

push|pull_request: develop
    |
    +--> server-develop

push: develop with server API changes
    |
    +--> server-swagger
```

## Migration Matrix

| Previous File | Disposition | Root Target | Trigger Strategy |
| --- | --- | --- | --- |
| `web/.github/PULL_REQUEST_TEMPLATE.md` | migrated | `.github/PULL_REQUEST_TEMPLATE.md` | GitHub contributor surface |
| `web/.github/ISSUE_TEMPLATE/1_bug_report.yml` | migrated | `.github/ISSUE_TEMPLATE/1_bug_report.yml` | GitHub contributor surface |
| `web/.github/ISSUE_TEMPLATE/2_feature_request.yml` | migrated | `.github/ISSUE_TEMPLATE/2_feature_request.yml` | GitHub contributor surface |
| `web/.github/ISSUE_TEMPLATE/3_question.yml` | migrated | `.github/ISSUE_TEMPLATE/3_question.yml` | GitHub contributor surface |
| `web/.github/ISSUE_TEMPLATE/4_other.md` | migrated | `.github/ISSUE_TEMPLATE/4_other.md` | GitHub contributor surface |
| `web/.github/workflows/issue-check-inactive.yml` | migrated | `.github/workflows/issue-check-inactive.yml` | `schedule` |
| `web/.github/workflows/issue-close-require.yml` | migrated | `.github/workflows/issue-close-require.yml` | `schedule` |
| `web/.github/workflows/issue-remove-inactive.yml` | migrated | `.github/workflows/issue-remove-inactive.yml` | `issues`, `issue_comment` |
| `web/.github/workflows/auto-merge.yml` | rewritten | `.github/workflows/auto-merge.yml` | `pull_request_target`, metadata-only |
| `web/.github/workflows/release.yml` | rewritten | `.github/workflows/web-release.yml` | `push` on `main|next|beta`, path-aware |
| `web/.github/workflows/publish-release-assets.yml` | rewritten | `.github/workflows/web-publish-release-assets.yml` | `release: published`, tag-gated on `web-v*` |
| `server/.github/workflows/develop.yaml` | rewritten | `.github/workflows/server-develop.yml` | `push` and `pull_request` on `develop`, path-aware |
| `server/.github/workflows/swagger.yaml` | rewritten | `.github/workflows/server-swagger.yml` | `push` on `develop`, path-aware |
| `server/.github/workflows/release.yml` | rewritten | `.github/workflows/server-release.yml` | `push` on `v*` tags |
| `.github/workflows/monorepo-check.yml` | rewritten | `.github/workflows/monorepo-check.yml` | `pull_request`, `push: main`, path-aware jobs |

## Trigger Notes

- `pull_request` and `push` validation jobs use changed-file detection to keep PR runs smaller while preserving full-sweep protection on `main`.
- `schedule`, `release`, `workflow_dispatch`, and `pull_request_target` cannot rely on path filters alone. These workflows use explicit event gating, tag naming, or metadata-only logic instead.
- Web release tags are namespaced as `web-v*`. Server release tags remain `v*`.

## Validation Contract

These files must remain true after future edits:

- Root validation is anchored on `Monorepo Summary`.
- No active workflow files are allowed under `web/.github/workflows/` or `server/.github/workflows/`.
- Web release automation must not wake server release automation.
- Publishing workflows must support no-op or shadow validation before any future cutover.
