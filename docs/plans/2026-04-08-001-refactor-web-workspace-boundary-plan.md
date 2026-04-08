---
title: "refactor: Web 工作区边界回收与根级编排收束"
type: refactor
status: completed
date: 2026-04-08
origin: external design doc "admin-feat-monorepo-baseline-design-20260408-073109.md"
---

# refactor: Web 工作区边界回收与根级编排收束

## Overview

把当前摊在仓库根目录的前端工作区整体回收到 `web/`，恢复 `server/` 与 `web/` 两个并列工程的顶层语义。根目录继续保留为跨工程编排层，但不再承担 Bun（Bun）/ Turborepo（Turbo）工作区根身份。

这次改动不重写运行时产品行为，也不顺手重做 split web 分发模型。目标很具体，文件系统边界追上已经存在的命令、构建、OpenAPI（OpenAPI）和发布合同。

## Problem Frame

当前仓库根目录直接持有：

- `apps/`
- `packages/`
- `package.json`
- `bun.lock`
- `turbo.json`
- `biome.json`
- `tsconfig.json`

这让仓库第一眼像一个纯前端 monorepo，但真实情况并不是这样。`server/` 本身已经是一个完整、独立、可单独理解的 Go 工程，拥有自己的 `go.mod`、CLI、路由、服务层、嵌入式静态资源目录和 Docker 构建方式。

结果就是根目录在同时说两套真相：

- 对贡献者，顶层像“前端 workspace + 一个 server 子目录”
- 对维护者，根 `Makefile`、根 `Dockerfile`、根 CI、根 OpenAPI 治理又在把仓库当成“两个并排工程 + 一个总控台”

这次计划要做的，就是把目录边界、命令入口、构建路径和文档说明统一成同一个说法。

## Requirements Trace

- R1. 顶层目录恢复 `server/` 与 `web/` 双工程语义，根目录不再直接保留 `apps/`、`packages/`、`bun.lock`、`turbo.json`、`biome.json`
- R2. `web/` 成为独立前端工作区，接管前端 lockfile、workspace 配置、前端脚本与依赖安装
- R3. 根 `package.json` 保留为极薄 orchestration shell，不再声明 workspace
- R4. 根 `Makefile` 保留 repo 级命令合同：`make bootstrap`、`make lint`、`make test`、`make format`、`make typecheck`、`make clean`、`make dev APP=...`、`make embed`、`make build-all`
- R5. 根级 `bun run openapi` 继续作为 OpenAPI 权威入口，`docs/openapi/*.json` 仍保留在根目录
- R6. canonical 构建链保持为嵌入式单镜像路径：根 `Dockerfile`、根 `docker-compose.yml`、`make embed`、`make build-all`
- R7. `server/web/admin-dist` 与 `server/web/user-dist` 继续作为嵌入式产物目标目录
- R8. 所有前端相对路径配置在迁移后仍成立，包括 OpenAPI 生成、Next.js（Next.js）构建、TypeScript（TypeScript）路径、Tailwind（Tailwind）/ PostCSS（PostCSS）配置和编辑器配置
- R9. `server/Dockerfile`、`docker/web.Dockerfile` 与 `.github/workflows/docker-build.yml` 必须完成兼容性审计，不能静默失配
- R10. README、AGENTS、API 治理文档、CI 与 Docker 说明更新到新坐标
- R11. 不引入新的运行时功能改动，不重做认证、会话、SSO、退出体验

## What Already Exists

已经存在的合同和实现不少，这次主要是复用和回收，不是重造：

- 根 `Makefile` 已经统一了全栈开发、嵌入式构建和清理入口
- 根 `package.json` 已经承载 `bun run openapi`、lint、typecheck、dev 等 repo 级入口
- 根 `Dockerfile` 已经是 all-in-one 嵌入式构建路径
- 根 `docker-compose.yml` 已经是默认单镜像 quickstart
- `server/web/admin-dist` 和 `server/web/user-dist` 已经存在，且被 `make embed` / `make build-all` 使用
- `.github/workflows/openapi-governance.yml` 已经把根 `bun run openapi` 当成权威治理入口
- `.github/workflows/docker-build.yml` 已经在发布 `server`、`admin-web`、`user-web`
- `redocly.yaml` 与 `docs/openapi/*.json` 已经是根级契约面
- `lefthook.yml` 已经是根级 hooks，总体方向正确，只是路径还停在旧根工作区

## NOT in Scope

- 不重设计 server API、错误响应格式或任何业务语义
- 不重做 `admin` / `user` 的认证、会话、SSO、退出或权限 UX
- 不把根级脚本体系重做成全新的产品形态，只做显式委托与边界收束
- 不把 split web 镜像产物模型从 `.next/standalone` 改成全新发布方案，本轮只做兼容性审计和最小修正
- 不引入新的前端测试框架或额外部署平台

## Step 0: Scope Challenge

### 已有能力复用

- 复用根 `Makefile` 作为总控台，而不是新建第二套顶层命令入口
- 复用根 `package.json` 作为极薄命令壳，而不是删掉后逼所有调用点知道 `web/`
- 复用 `docs/openapi/` 作为 server 与 web 共享契约面，而不是把 spec 搬进 `web/`
- 复用根 `Dockerfile` 和根 `docker-compose.yml` 作为 canonical 发布链

### 最小变更集

这次只做能达成目标的最小闭环：

1. 把前端工作区整体迁到 `web/`
2. 收窄根目录到 orchestration shell
3. 修正所有相对路径与构建上下文
4. 保住 canonical 嵌入式链路
5. 审计兼容链，不做额外产品重设计

### Complexity Check

这次一定会触碰 8 个以上文件，这是目录边界回收的天然 blast radius，不是过度设计本身。真正会把事情做炸的，不是文件数，而是把以下两件本不属于本轮的事一起拖进来：

- split web 分发模型重写
- 认证 / 会话产品语义重做

这两个都已经被剥离出本轮，所以 scope 已经被收束到合理边界。

### Search Check

本轮不引入新框架或新基础设施，主要是重新安放已有 Bun、Turbo、Next、Docker、Redocly（Redocly）合同。外部搜索不是关键，repo 内既有实现和前序 learnings 更重要。

已直接应用的既有 learnings：

- `biome-config-path-from-root`，根目录调用 Biome（Biome）时必须显式指向 `web/biome.json`
- `turbo-format-no-workspace-script`，格式化不能依赖 `turbo run format` 这种会静默漏跑的路径
- `monorepo-glue-layer-first`，保留极薄根壳比一步删干净更安全

### TODO Cross-reference

仓库当前没有 `TODOS.md`，本计划不依赖已有 deferred work。

### Completeness Check

这轮按完整版本执行，不走“只搬目录，不修入口”的捷径。目录迁移如果不同时修正 Docker、CI、OpenAPI 和文档，只会得到一个看起来干净、第一次运行就撞墙的假完成。

### Distribution Check

本轮涉及的 artifact 只有既有三条链：

```text
Canonical path in this phase
============================
web/apps/* --build/export--> server/web/*-dist
                         \--> root Dockerfile
                         \--> docker-compose.yml
                         \--> make embed / make build-all

Compatibility paths to audit
============================
server/Dockerfile -----------------> server image
docker/web.Dockerfile -------------> admin-web / user-web images
```

结论，canonical 链纳入实现与验证，compatibility 链纳入审计与最小修正。

## Context & Research

### Relevant Files

- 根入口：`Makefile`、`package.json`、`Dockerfile`、`docker-compose.yml`、`lefthook.yml`
- OpenAPI：`docs/openapi/*.json`、`redocly.yaml`、`docs/api-governance.md`、`.github/workflows/openapi-governance.yml`
- 前端工作区：`apps/*`、`packages/*`、`turbo.json`、`biome.json`、`tsconfig.json`、`scripts/*`
- 构建链：`server/Dockerfile`、`docker/web.Dockerfile`、`.github/workflows/docker-build.yml`
- 编辑器与文档：`.vscode/settings.json`、`README.md`、`AGENTS.md`

### Current Pitfalls Confirmed In Code

- 根 `Dockerfile` 直接复制根 `apps/` 与 `packages/`，迁移后一定失效
- 根 `Makefile` 直接调用 `bun run ...`，迁移后必须显式委托到 `web/`
- `.dockerignore` 当前忽略的是根 `apps/*/.next` 与 `.turbo`，迁移后必须改成 `web/...`
- `apps/*/openapi-ts.config.ts` 现在使用 `../../docs/openapi/*`，迁入 `web/` 后层级会多一层
- `apps/*/next.config.ts` 中 `monaco-themes/themes` alias 现在指向 `../../node_modules/...`，迁移后路径会变化
- `docker/web.Dockerfile` 仍按 `.next/standalone` 打包，与当前 `output: "export"` 冲突，这就是 compatibility 链必须显式标注的原因

## Key Technical Decisions

### D1. 保留极薄根 `package.json`

根 `package.json` 保留，但职责缩到三件事：

- repo 级命令入口
- OpenAPI 治理入口
- hooks / CI 兼容壳

不再声明 `workspaces`，也不再承担实际依赖根。

### D2. `web/` 拥有全部前端工具链状态

`web/` 将接管：

- `apps/`
- `packages/`
- `bun.lock`
- `turbo.json`
- `biome.json`
- `tsconfig.json`
- 前端专用 `scripts/`

### D3. `docs/openapi/` 保持根级共享契约

OpenAPI 导出与生成链路继续是：

```text
server/cmd/openapi.go
    │
    ├──> docs/openapi/admin.json
    ├──> docs/openapi/common.json
    └──> docs/openapi/user.json
             │
             └──> web/apps/*/openapi-ts*.config.ts
                      │
                      └──> web/apps/*/services/*-api
```

### D4. canonical 与 compatibility 明确拆开

- canonical，本轮必须可用：
  - 根 `Dockerfile`
  - 根 `docker-compose.yml`
  - `make embed`
  - `make build-all`
- compatibility，本轮必须被审计：
  - `server/Dockerfile`
  - `docker/web.Dockerfile`
  - `.github/workflows/docker-build.yml`

### D5. 根目录仍是总控台，但所有前端命令显式进入 `web/`

委托模式统一成显式路径，不再赌当前工作目录：

```text
root command
   │
   ├── make bootstrap -----> cd web && bun install
   ├── make lint ----------> cd web && bun run lint
   ├── make typecheck -----> cd web && bun run typecheck
   ├── make embed ---------> cd web && bun run build --filter=...
   └── bun run openapi ----> cd web && bun run openapi:client
```

## Open Questions

### Resolved

- 是否删除根 `package.json`，不删，保留极薄 orchestration shell
- canonical 发布链是哪条，嵌入式单镜像路径
- `docs/openapi/` 是否搬进 `web/`，不搬，继续保留根级共享契约

### Deferred

- split web 镜像是否继续作为正式支持面
- `docker/web.Dockerfile` 是否未来切换到静态站点镜像模型
- 是否额外保留 `make build` 作为与 `make build-all` 同等强度的发布验收门槛

## Implementation Units

- [x] **Unit 1: 建立 `web/` 工作区根并迁移前端目录**

**Goal:** 把前端源码、配置、lockfile 和前端脚本整体搬入 `web/`

**Files / Paths:**
- Move: `apps/` → `web/apps/`
- Move: `packages/` → `web/packages/`
- Move: `bun.lock` → `web/bun.lock`
- Move: `turbo.json` → `web/turbo.json`
- Move: `biome.json` → `web/biome.json`
- Move: `tsconfig.json` → `web/tsconfig.json`
- Move: `scripts/` → `web/scripts/`
- Create: `web/package.json`
- Modify: 根 `package.json`

**Notes:**
- 根 `package.json` 只保留 repo 命令，不再保留 `workspaces`
- `web/package.json` 承接原先工作区依赖、workspace 声明和前端 scripts

**Verification:**
- `cd web && bun install`
- `cd web && bun run build`
- `cd web && bun run typecheck`

---

- [x] **Unit 2: 收窄根级 orchestration shell 与 hooks**

**Goal:** 保住 repo 级入口，同时把所有前端命令改成显式委托

**Files:**
- Modify: `package.json`
- Modify: `Makefile`
- Modify: `lefthook.yml`

**Notes:**
- 所有根级 Bun 命令使用 `cd web && bun ...`
- `make dev`、`make embed`、`make build-all`、`make clean` 都改成新坐标
- hooks 如果继续从根触发 Biome / commitlint，必须显式指向 `web/`

**Verification:**
- `make bootstrap`
- `make lint`
- `make format`
- `make typecheck`
- `make dev APP=admin`
- `make dev APP=user`

---

- [x] **Unit 3: 修正 workspace-relative 配置与编辑器坐标**

**Goal:** 让 `web/` 内部构建和本地编辑体验恢复正常

**Files:**
- Modify: `web/apps/*/openapi-ts*.config.ts`
- Modify: `web/apps/*/next.config.ts`
- Modify: `web/apps/*/tsconfig.json`
- Modify: `web/apps/*/components.json`
- Modify: `web/apps/*/postcss.config.*`
- Modify: `.vscode/settings.json`

**Notes:**
- OpenAPI spec 路径相对层级整体加一层
- `monaco-themes/themes` alias 指向的 `node_modules` 层级整体加一层
- Tailwind / IDE 配置路径从根 `apps/*`、`packages/*` 改成 `web/...`

**Verification:**
- `cd web && bun run dev:admin`
- `cd web && bun run dev:user`
- `cd web && bun run openapi`

---

- [x] **Unit 4: 修正 canonical 嵌入式构建链**

**Goal:** 让根 `Dockerfile`、根 `docker-compose.yml`、`make embed`、`make build-all` 继续成立

**Files:**
- Modify: `Dockerfile`
- Modify: `docker-compose.yml`
- Modify: `.dockerignore`
- Modify: `Makefile`

**Notes:**
- 根 `Dockerfile` 必须复制 `web/...`，并在 `/app/web` 执行 Bun 构建
- `.dockerignore` 必须改成忽略 `web/apps/*/.next`、`web/.turbo`、`web/node_modules`
- `make embed` 的复制源改成 `web/apps/*/out`

**Verification:**
- `make embed`
- `make build-all`
- `docker build -f Dockerfile .`
- `docker compose build`

---

- [x] **Unit 5: 审计并修正 compatibility 链**

**Goal:** 明确 `server` / `admin-web` / `user-web` 镜像链当前是否仍成立

**Files:**
- Modify: `server/Dockerfile`
- Modify: `docker/web.Dockerfile`
- Modify: `.github/workflows/docker-build.yml`

**Notes:**
- `server/Dockerfile` 至少要确认上下文假设在新目录下不失效
- `docker/web.Dockerfile` 如果继续保留，必须先修路径，再明确它仍然是 compatibility path
- 如果某条链本轮无法做真，需要在文档里显式 deferred，而不是假装它可用

**Verification:**
- `docker build -f server/Dockerfile server`
- `docker build --build-arg APP_NAME=admin -f docker/web.Dockerfile .`
- `docker build --build-arg APP_NAME=user -f docker/web.Dockerfile .`

---

- [x] **Unit 6: 文档、治理与映射表回写**

**Goal:** 让说明文字与仓库现实同步

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/api-governance.md`
- Modify: `redocly.yaml`
- Modify: `.github/workflows/openapi-governance.yml`

**Notes:**
- 把根目录描述成总控台，不再描述成前端 workspace 根
- 在 README 中补 `旧路径 -> 新路径`、`旧命令 -> 新命令` 映射

**Verification:**
- `bun run openapi`
- README 首屏能解释清楚 `server/`、`web/`、根目录三者职责

## Test Review

这类改动的风险不在业务逻辑，而在“命令入口、相对路径和构建上下文 quietly broken”。所以测试重点是命令烟雾（smoke）和构建烟雾，而不是额外业务单测。

```text
CODE PATH COVERAGE
===========================
[+] Root command contract
    │
    ├── [PASS] make bootstrap
    ├── [BLOCKED] make lint — existing `unused` finding in `server/worker/task/deps.go`
    ├── [BLOCKED] make format — local `goimports` missing
    ├── [PASS] make typecheck
    ├── [DEFERRED] make dev APP=admin|user — not required for this filesystem/build boundary pass
    └── [PASS] bun run openapi

[+] Embedded build path
    │
    ├── [PASS] make embed
    ├── [PASS] make build-all
    ├── [PASS] docker build -f Dockerfile .
    └── [PASS] docker compose build

[+] Compatibility image path
    │
    ├── [PASS] docker build -f server/Dockerfile server
    ├── [PASS] docker build -f docker/web.Dockerfile --build-arg APP_NAME=admin .
    └── [PASS] docker build -f docker/web.Dockerfile --build-arg APP_NAME=user .

[+] Config path rewrites
    │
    ├── [PASS] OpenAPI generator relative paths
    ├── [PASS] Next.js webpack alias paths
    └── [PASS] VS Code Tailwind config paths

─────────────────────────────────
COVERAGE: 13 pass / 2 blocked / 1 deferred across 16 tracked paths
QUALITY:  root command + build smoke mostly green; remaining gaps are explicit
GAPS: `make lint` 受既有 lint 问题阻塞，`make format` 受本机缺 `goimports` 阻塞，`make dev` 本轮 deferred
─────────────────────────────────
```

### Required Verification Matrix

- V1. `cd web && bun install`
- V2. `cd web && bun run build`
- V3. `cd web && bun run typecheck`
- V4. `make bootstrap`
- V5. `make lint`
- V6. `make format`
- V7. `make typecheck`
- V8. `make embed`
- V9. `make build-all`
- V10. `bun run openapi`
- V11. `docker build -f Dockerfile .`
- V12. `docker compose build`
- V13. compatibility 链最小构建验证

## Failure Modes

| Codepath | Realistic failure | Test / smoke | Error handling | User-visible result | Critical gap |
|---|---|---|---|---|---|
| 根 `package.json` 委托 | 根命令仍引用旧根 workspace，`bun run` 直接失败 | V4-V7, V10 | shell 非零退出 | 开发者第一步就撞墙 | 是 |
| OpenAPI 生成 | `openapi-ts` 相对路径少一层，client 生成失败 | V10 | CLI 非零退出 | CI 红，SDK 不更新 | 是 |
| Next 构建 alias | `monaco-themes/themes` 指向旧 `node_modules` | V2, V8 | build 非零退出 | admin/user 构建失败 | 是 |
| `make embed` | 复制源仍是根 `apps/*/out`，嵌入目录为空 | V8-V9 | shell 非零退出或复制空目录 | 产物缺静态资源 | 是 |
| 根 `Dockerfile` | Docker context 仍复制根 `apps/`/`packages/` | V11 | build 非零退出 | 默认部署链失效 | 是 |
| `.dockerignore` | 忽略规则未更新，构建上下文变大或遗漏关键文件 | V11-V12 | 无显式处理 | 构建慢或失败 | 否 |
| `docker/web.Dockerfile` | 继续假设 `.next/standalone`，与 `output: "export"` 冲突 | V13 | build 非零退出 | compatibility 镜像不可用 | 否，本轮显式审计 |
| VS Code 配置 | Tailwind 配置路径仍指向根 `apps/*` | 静态 inspection | 无 | 编辑器补全失效 | 否 |

本计划里的 critical gaps 都有对应 smoke/build 验证，不能跳过。

## System-Wide Impact

- 文件系统边界改变，但运行时产品行为不变
- 根级命令合同保留，调用实现从“直接运行”改成“显式委托”
- canonical 发布链继续是嵌入式单镜像部署
- compatibility 链被从“默认认为可用”降级成“显式审计并标注状态”

## Worktree Parallelization Strategy

| Step | Modules touched | Depends on |
|---|---|---|
| Unit 1 | root, `web/`, `apps/`, `packages/`, `scripts/` | — |
| Unit 2 | root, `web/`, hooks | Unit 1 |
| Unit 3 | `web/apps/`, `.vscode/` | Unit 1 |
| Unit 4 | root build assets, `server/web/` | Units 1-3 |
| Unit 5 | `server/`, `docker/`, `.github/workflows/` | Unit 4 |
| Unit 6 | docs, root metadata | Units 1-5 |

Parallel lanes:

- Lane A: Unit 1 → Unit 2 → Unit 4
- Lane B: Unit 3
- Lane C: Unit 5
- Lane D: Unit 6

Execution order:

1. 先做 Lane A，建立真实目录和根级委托关系
2. 同时开 Lane B，修正 `web/apps/*` 内部相对路径
3. A + B 合并后做 Lane C，审计 compatibility 构建链
4. 最后做 Lane D，统一文档和治理说明

Conflict flags:

- Lane A 与 Lane B 都会碰 `web/`，如果并行开发，必须按目录拆 ownership
- Lane C 会读取 Lane A 的新路径，所以不能先于 Lane A 合并

## Progress Update

2026-04-08 当前实现进度：

- Unit 1-6 已落地，`web/` 工作区、根级 orchestration shell、OpenAPI 路径、canonical 构建链、compatibility 构建链与文档说明都已按新边界收口
- canonical 验证已通过：`cd web && bun install`、`cd web && bun run build`、`cd web && bun run typecheck`、`make bootstrap`、`make typecheck`、`make embed`、`make build-all`、`bun run openapi`、`docker build -f Dockerfile .`、`docker compose build`
- compatibility 验证已通过：`docker build -f server/Dockerfile server`、`docker build --build-arg APP_NAME=admin -f docker/web.Dockerfile .`、`docker build --build-arg APP_NAME=user -f docker/web.Dockerfile .`
- 额外 smoke：`make test`、`make clean` 已通过
- 已知阻塞项有 2 个：
  - `make lint` 仍会因既有问题失败：`server/worker/task/deps.go` 的 `Deps.currentConfig` 未使用
  - `make format` 在当前机器上因缺少 `goimports` 失败，这属于环境依赖缺失，不是本轮 `web/` 边界改动引入的问题
- `cd web && bun run build` 在静态导出阶段仍会打印 `ENVIRONMENT_FALLBACK` 警告，但构建最终退出码为 0

## Verification Checklist

- [x] 根目录不再存在 `apps/`、`packages/`、`bun.lock`、`turbo.json`、`biome.json`
- [x] `web/` 内存在完整前端工作区根
- [x] 根 `package.json` 不再声明 `workspaces`
- [x] `cd web && bun install` 成功
- [x] `cd web && bun run build` 成功
- [x] `cd web && bun run typecheck` 成功
- [x] `make bootstrap` 成功
- [ ] `make lint` 成功
- [ ] `make format` 成功
- [x] `make typecheck` 成功
- [x] `make embed` 成功
- [x] `make build-all` 成功
- [x] `bun run openapi` 成功
- [x] `docker build -f Dockerfile .` 成功
- [x] `docker compose build` 成功
- [x] compatibility 链构建状态被验证并记录

## Sources & References

- 2026-04-08 外部设计稿：`admin-feat-monorepo-baseline-design-20260408-073109.md`
- `docs/brainstorms/2026-04-03-monorepo-dx-unification-requirements.md`
- `docs/plans/2026-04-03-001-refactor-monorepo-dx-unification-plan.md`
- `Makefile`
- `package.json`
- `Dockerfile`
- `docker-compose.yml`
- `server/Dockerfile`
- `docker/web.Dockerfile`
- `.github/workflows/openapi-governance.yml`
- `.github/workflows/docker-build.yml`
- `.dockerignore`
- `redocly.yaml`
- `docs/api-governance.md`

## GSTACK REVIEW REPORT

| Review | Trigger | Why | Runs | Status | Findings |
|--------|---------|-----|------|--------|----------|
| CEO Review | `/plan-ceo-review` | Scope & strategy | 1 | CLEAR | scope 已收束到仓库边界、构建发布边界与契约真相 |
| Codex Review | `/codex review` | Independent 2nd opinion | 0 | — | — |
| Eng Review | `/plan-eng-review` | Architecture & tests (required) | 1 | CLEAR | canonical/compatibility 拆分明确，13 个关键验证路径已写入计划 |
| Design Review | `/plan-design-review` | UI/UX gaps | 0 | — | — |
| DX Review | `/plan-devex-review` | Developer experience gaps | 0 | — | — |

**UNRESOLVED:** 0
**VERDICT:** CEO + ENG CLEARED — implementation started, core verification mostly green

## Closeout Note

收口映射如下：

- `README.md` / `AGENTS.md` 固定 support matrix 与版本真相
- `.github/workflows/monorepo-boundary-guardrail.yml` 固定 boundary guardrail
- `.github/workflows/repo-contracts.yml` 固定 `bun run repo:contracts`（根命令合同 + OpenAPI lint）与 `docker build -f Dockerfile .`（canonical image smoke）这两道独立 gate

compatibility lane 和少量 deferred 项目前仍保留人工审计尾巴，不把它们伪装成 100% 自动化闭环。
