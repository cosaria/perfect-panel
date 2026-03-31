# Web Biome Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace ESLint with Biome across the web workspace so `lint` only runs Biome, while type checking moves to dedicated `check-types` and root `typecheck` commands.

**Architecture:** Update package scripts so every web package exposes a Biome-only `lint` and a separate TypeScript `check-types` where applicable. Remove the shared ESLint config package and ESLint config files, then align root tooling and CI to the new split between linting and type checking.

**Tech Stack:** Bun workspaces, Turbo, Biome, TypeScript, Next.js, GitHub Actions

---

### Task 1: Rewire Workspace Scripts

**Files:**
- Modify: `web/package.json`
- Modify: `web/apps/user/package.json`
- Modify: `web/apps/admin/package.json`
- Modify: `web/packages/ui/package.json`

- [ ] Replace package-level `lint` commands so they run Biome only.
- [ ] Add `check-types` commands to packages that own TypeScript source.
- [ ] Add a root `typecheck` script that runs `turbo check-types`.
- [ ] Update any root helper scripts that still invoke ESLint directly.

### Task 2: Reconfigure Biome Scope

**Files:**
- Modify: `web/biome.json`

- [ ] Preserve generated-code exclusions that were previously handled by ESLint ignores, especially the `services/**/*.ts` API clients under both apps.
- [ ] Keep Biome focused on source files instead of generated Next output noise.

### Task 3: Remove ESLint Tooling

**Files:**
- Delete: `web/eslint.config.js`
- Delete: `web/apps/user/eslint.config.js`
- Delete: `web/apps/admin/eslint.config.js`
- Delete: `web/packages/ui/eslint.config.js`
- Delete: `web/packages/eslint-config/package.json`
- Delete: `web/packages/eslint-config/base.js`
- Delete: `web/packages/eslint-config/next.js`
- Delete: `web/packages/eslint-config/react-internal.js`

- [ ] Remove `eslint` and `@workspace/eslint-config` dependencies from workspace package manifests.
- [ ] Delete ESLint config entrypoints and the shared ESLint config workspace package.

### Task 4: Align CI Entrypoints

**Files:**
- Modify: `.github/workflows/monorepo-check.yml`

- [ ] Keep CI invoking `bun run lint`, but add `bun run check-types` before `bun run build` for both web apps.

### Task 5: Refresh Dependency Lock And Verify

**Files:**
- Modify: `web/bun.lock`

- [ ] Refresh the Bun lockfile after dependency removal.
- [ ] Run `bun run lint` in `web/apps/user` and `web/apps/admin`.
- [ ] Run `bun run check-types` in `web/apps/user`, `web/apps/admin`, and `web/packages/ui`.
- [ ] Run `bun run build` in both apps to confirm the new split does not hide compile issues.