---
title: "refactor: monorepo baseline 收尾与仓库真相封口"
type: refactor
status: completed
date: 2026-04-08
origin: docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md
---

# Monorepo Baseline 收尾 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: 使用 `subagent-driven-development`（推荐）或 `executing-plans` 按任务逐项执行。步骤使用 checkbox（`- [ ]`）语法追踪。

**Goal:** 把当前已经基本成型的 `server/` + `web/` 双工程仓库，收尾成一套“文档、命令、CI、Docker、发布支持面都说同一句话”的正式基线。

**Architecture:** 这次不再重做目录，不再碰运行时产品行为，也不再扩 scope。收尾工作的核心是四件事：锁定官方支持面、补机器护栏、统一工具链真相、把发布验证升级为长期 gate。默认保留根 `Dockerfile` + `make embed` + `make build-all` 作为 canonical 路径；`server`、`admin-web`、`user-web` 镜像链先按“兼容路径”处理，只有在 smoke gate 跑通后才升级为正式支持面。

**Tech Stack:** Go 1.25、Bun 1.3.0、Turbo、Next.js 15、GitHub Actions、Docker、Redocly、lefthook

---

## Scope Check

这份计划只处理 monorepo baseline 的“封口”问题，不再继续扩展到认证、会话、SSO、前端交互或运行时架构。它是对 [2026-04-08-001-refactor-web-workspace-boundary-plan.md](/Users/admin/Codes/ProxyCode/perfect-panel/docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md) 的 closeout plan，不是新一轮仓库重构。

## File Structure

在拆任务前，先锁定这轮收尾会触达的文件和职责边界。

### 新建文件

- `.github/scripts/check-monorepo-boundary.sh`
  - 负责检查仓库根目录是否重新出现前端 workspace 根文件或目录
  - 失败时输出明确文件名和修复建议
- `.github/workflows/monorepo-boundary-guardrail.yml`
  - 负责在 PR / push 时执行 boundary guardrail
- `.github/workflows/repo-contracts.yml`
  - 负责跑 canonical release gate 和文档/契约 smoke

### 修改文件

- `README.md`
  - 对外声明仓库真相、支持矩阵、最低版本要求、canonical / compatibility 路径
- `AGENTS.md`
  - 同步命令、依赖版本、仓库架构描述和发布真相
- `server/README.md`
  - 明确它是独立 Go 工程，但默认官方发布路径由仓库根 orchestrate
- `package.json`
  - 保持极薄 shell 身份，不新增 workspace 语义，只补 release-contract 级命令
- `Makefile`
  - 清理旧根目录残影，补充收尾验证入口
- `.github/workflows/openapi-governance.yml`
  - 保持 OpenAPI 为根级合同，避免和新 repo-contracts 重复或漂移
- `.github/workflows/docker-build.yml`
  - 给 compatibility 镜像链加注释、条件或支持级别说明
- `Dockerfile`
  - 如有必要，只补注释或显式验证，不改变 canonical 构建语义
- `server/Dockerfile`
  - 如有必要，只补用途说明，不引入与根发布链冲突的新行为
- `docker/web.Dockerfile`
  - 如有必要，只补“静态导出镜像”语义和 smoke 约束，不改产品行为
- `docs/api-governance.md`
  - 记录 OpenAPI 仍是根级共享契约，说明与 `web/` 的边界
- `docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md`
  - 在收尾完成后改为 `completed`，追加 closeout note

## Decision Record

这份计划先锁定 3 个决策，后续执行不再来回摇摆：

1. **Canonical 路径仍是根级嵌入式链路**
   - 根 `Dockerfile`
   - 根 `docker-compose.yml`
   - `make embed`
   - `make build-all`
2. **Compatibility 路径先按“兼容但非默认”处理**
   - `server/Dockerfile`
   - `docker/web.Dockerfile`
   - `.github/workflows/docker-build.yml` 里的 `server`、`admin-web`、`user-web`
3. **版本真相以实际工具链为准**
   - Go 最低版本按 `server/go.mod` 的 `go 1.25.0`
   - Bun 最低版本按 `package.json` / `web/package.json` 的 `bun@1.3.0`
   - 文档不再保留 `Go 1.21+` / `Go 1.23+` 这种多版本口径

## Verification Matrix

收尾完成后，至少要能稳定通过下面这些命令：

```bash
make typecheck
cd server && go test ./...
cd server && tmpdir=$(mktemp -d) && go run . openapi -o "$tmpdir"
bun run openapi:lint
docker build -f Dockerfile .
bash .github/scripts/check-monorepo-boundary.sh
```

Compatibility 审计命令：

```bash
docker build -f server/Dockerfile server
docker build -f docker/web.Dockerfile --build-arg APP_NAME=admin .
docker build -f docker/web.Dockerfile --build-arg APP_NAME=user .
```

## Task 1: 锁定支持面并把仓库真相写死

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `server/README.md`
- Modify: `.github/workflows/docker-build.yml`
- Modify: `docs/api-governance.md`

- [ ] **Step 1: 审计当前支持面文案**

Run:

```bash
rg -n 'Go 1\\.21\\+|Go 1\\.23\\+|compatibility path|canonical|admin-web|user-web|server image' \
  README.md AGENTS.md server/README.md docs/api-governance.md .github/workflows/docker-build.yml
```

Expected:

```text
能列出所有旧版本口径和支持面描述位置，方便一次性统一
```

- [ ] **Step 2: 在 `README.md` 新增“Support Matrix”段落**

写入类似下面的结构，明确默认路径和兼容路径：

```md
## Support Matrix

| Surface | Status | Purpose |
|---|---|---|
| Root `Dockerfile` + `make embed` + `make build-all` | supported | 默认官方发布链 |
| Root `docker-compose.yml` | supported | 默认本地/单镜像 quickstart |
| `server/Dockerfile` | compatibility | 独立 server 镜像链，非默认入口 |
| `docker/web.Dockerfile` (`admin-web` / `user-web`) | compatibility | 独立静态前端镜像链，非默认入口 |
```

- [ ] **Step 3: 把 `AGENTS.md` 中的版本和仓库真相统一到实际基线**

至少统一下面两处：

```md
- Go 1.25+、golangci-lint、goimports
- Node.js 20+、Bun 1.3.0+
```

并在架构说明里补一句：

```md
默认官方发布链是根 `Dockerfile` / `docker-compose.yml` / `make embed` / `make build-all`。
`server`、`admin-web`、`user-web` 镜像链为兼容路径，不是新贡献者默认入口。
```

- [ ] **Step 4: 给 `server/README.md` 加一段“发布边界说明”**

补充类似下面的句子，避免读者以为 `server/` 自己就是官方发布入口：

```md
`server/` 是独立 Go 工程，但仓库默认官方发布路径由根目录统一编排。
需要全栈嵌入式构建时，请优先使用仓库根的 `make embed` / `make build-all` / 根 `Dockerfile`。
```

- [ ] **Step 5: 给 `.github/workflows/docker-build.yml` 加用途注释**

在 job 级别补注释，明确：

```yaml
# server / admin-web / user-web are compatibility image lanes.
# Canonical release path remains the root all-in-one Dockerfile.
```

- [ ] **Step 6: 运行文档漂移检查**

Run:

```bash
rg -n 'Go 1\\.21\\+|Go 1\\.23\\+|Go 1\\.25\\+|compatibility|supported|canonical' \
  README.md AGENTS.md server/README.md docs/api-governance.md .github/workflows/docker-build.yml
```

Expected:

```text
只剩一套版本真相，且 canonical / compatibility 语义一致
```

- [ ] **Step 7: 提交这一批文档与声明收口**

```bash
git add README.md AGENTS.md server/README.md docs/api-governance.md .github/workflows/docker-build.yml
git commit -m "docs: lock monorepo support matrix and toolchain truth"
```

## Task 2: 增加根目录边界护栏，防止 workspace 回流

**Files:**
- Create: `.github/scripts/check-monorepo-boundary.sh`
- Create: `.github/workflows/monorepo-boundary-guardrail.yml`
- Modify: `README.md`
- Modify: `AGENTS.md`

- [ ] **Step 1: 写边界检查脚本**

创建 `.github/scripts/check-monorepo-boundary.sh`，核心逻辑直接写死，不做花哨抽象：

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT_FORBIDDEN=(
  "apps"
  "packages"
  "bun.lock"
  "turbo.json"
  "biome.json"
  "tsconfig.json"
)

violations=()

for path in "${ROOT_FORBIDDEN[@]}"; do
  if [ -e "$path" ]; then
    violations+=("$path")
  fi
done

if [ "${#violations[@]}" -gt 0 ]; then
  echo "ERROR: monorepo boundary regression detected."
  printf ' - %s\n' "${violations[@]}"
  echo "Move frontend workspace state back under web/."
  exit 1
fi

echo "OK: root workspace boundary is clean."
```

- [ ] **Step 2: 本地执行脚本，确认当前基线通过**

Run:

```bash
bash .github/scripts/check-monorepo-boundary.sh
```

Expected:

```text
OK: root workspace boundary is clean.
```

- [ ] **Step 3: 新增 guardrail workflow**

创建 `.github/workflows/monorepo-boundary-guardrail.yml`：

```yaml
name: Monorepo Boundary Guardrail

on:
  pull_request:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  boundary-guardrail:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Verify root/frontend boundary
        run: bash .github/scripts/check-monorepo-boundary.sh
```

- [ ] **Step 4: 在 `README.md` 和 `AGENTS.md` 补一句 guardrail 存在**

加入一句明确说明：

```md
仓库通过 `Monorepo Boundary Guardrail` 阻止前端 workspace 根文件重新回流到仓库根目录。
```

- [ ] **Step 5: 提交 boundary 护栏**

```bash
git add .github/scripts/check-monorepo-boundary.sh .github/workflows/monorepo-boundary-guardrail.yml README.md AGENTS.md
git commit -m "ci: add monorepo boundary guardrail"
```

## Task 3: 统一根级命令合同与版本真相

**Files:**
- Modify: `Makefile`
- Modify: `package.json`
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `server/README.md`

- [ ] **Step 1: 清理 `Makefile` 里的旧根目录残影**

把 [web-clean](/Users/admin/Codes/ProxyCode/perfect-panel/Makefile#L83) 从：

```make
rm -rf .turbo web/apps/*/.next web/apps/*/.turbo web/.turbo
```

改成：

```make
rm -rf web/apps/*/.next web/apps/*/.turbo web/.turbo
```

- [ ] **Step 2: 在根 `package.json` 增加一个 repo-contracts 聚合命令**

新增：

```json
"repo:contracts": "make typecheck && make test && bun run openapi:lint"
```

不要重新引入 workspace、不要把依赖装回根目录。

- [ ] **Step 3: 把 README / AGENTS / server README 里的最低版本统一为 Go 1.25+**

统一后应该满足：

```text
README.md      -> Go 1.25+
AGENTS.md      -> Go 1.25+
server/README.md -> Go 1.25+
```

- [ ] **Step 4: 跑根级命令合同 smoke**

Run:

```bash
make typecheck
cd server && go test ./...
bun run repo:contracts
```

Expected:

```text
全部命令通过，且不要求把前端 workspace 文件搬回根目录
```

- [ ] **Step 5: 提交命令合同与版本真相同步**

```bash
git add Makefile package.json README.md AGENTS.md server/README.md
git commit -m "chore: sync repo contracts and toolchain versions"
```

## Task 4: 把发布验证升级为长期 gate，并关闭 baseline 计划

**Files:**
- Create: `.github/workflows/repo-contracts.yml`
- Modify: `.github/workflows/openapi-governance.yml`
- Modify: `.github/workflows/docker-build.yml`
- Modify: `docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md`
- Modify: `README.md`

- [ ] **Step 1: 新增 repo-contracts workflow**

创建 `.github/workflows/repo-contracts.yml`：

```yaml
name: Repo Contracts

on:
  pull_request:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  repo-contracts:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: server/go.mod
      - uses: oven-sh/setup-bun@v2
        with:
          bun-version: 1.3.0
      - name: Install web dependencies
        run: cd web && bun install --frozen-lockfile
      - name: Run repo contracts
        run: bun run repo:contracts
      - name: Build canonical image
        run: docker build -f Dockerfile .
```

- [ ] **Step 2: 明确 compatibility audit 的执行方式**

在 `repo-contracts.yml` 或 `docker-build.yml` 里二选一，选一个固定下来：

```yaml
- name: Build compatibility server image
  run: docker build -f server/Dockerfile server

- name: Build compatibility admin-web image
  run: docker build -f docker/web.Dockerfile --build-arg APP_NAME=admin .

- name: Build compatibility user-web image
  run: docker build -f docker/web.Dockerfile --build-arg APP_NAME=user .
```

推荐先放在 `workflow_dispatch` 或单独 matrix job 下，避免第一次把 PR 时长炸穿。

- [ ] **Step 3: 保持 `openapi-governance.yml` 只负责契约治理**

不要把 Docker build、Go 全量测试和 boundary guardrail 再塞进去。它应该继续只做：

```yaml
- bun install
- bun run openapi
```

目的很简单，避免一个 workflow 同时承载三种职责，之后没人知道到底哪条合同挂了。

- [ ] **Step 4: 关闭主 baseline 计划**

在 `docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md` 执行两处修改：

1. frontmatter:

```yaml
status: completed
```

2. 末尾追加 closeout note：

```md
## Closeout Note

Monorepo baseline is considered closed when:

- support matrix is documented
- root/frontend boundary guardrail is enforced by CI
- toolchain version truth is unified to Go 1.25+ / Bun 1.3.0+
- repo-contract smoke and canonical Docker build run in CI
```

- [ ] **Step 5: 跑最终 smoke**

Run:

```bash
bash .github/scripts/check-monorepo-boundary.sh
bun run repo:contracts
docker build -f Dockerfile .
```

如果 compatibility 路径这轮也纳入 gate，再补跑：

```bash
docker build -f server/Dockerfile server
docker build -f docker/web.Dockerfile --build-arg APP_NAME=admin .
docker build -f docker/web.Dockerfile --build-arg APP_NAME=user .
```

- [ ] **Step 6: 提交 closeout**

```bash
git add .github/workflows/repo-contracts.yml .github/workflows/openapi-governance.yml .github/workflows/docker-build.yml docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md README.md
git commit -m "ci: close out monorepo baseline contracts"
```

## Rollback Plan

如果收尾过程中某一步把 CI 或构建链搞坏，按下面顺序回退：

1. 优先 `git revert` 当前 task 对应的单个 commit
2. 不回退 `web/` 目录迁移本身，不把前端 workspace 状态搬回根目录
3. 若新 workflow 导致 PR 卡死，先在同一分支移除 workflow 或改成 `workflow_dispatch`，不要顺手改 Docker / OpenAPI 逻辑
4. 若 compatibility 镜像 smoke 失败，允许把它保留为 manual audit，不要影响 canonical 链的完成状态

## Exit Criteria

这份计划完成时，必须同时满足下面 8 条：

1. `README.md`、`AGENTS.md`、`server/README.md` 对工具链最低版本说同一句话
2. `README.md` 明确存在 canonical / compatibility support matrix
3. 根目录重新出现 `apps/`、`packages/`、`bun.lock`、`turbo.json`、`biome.json` 会被 CI 直接拦下
4. 根 `package.json` 仍是极薄 orchestration shell，没有重新声明 workspace
5. `make typecheck`、`go test ./...`、`bun run openapi:lint` 可以作为长期 repo contract 跑通
6. 根 `Dockerfile` 构建进入 CI gate
7. `docs/plans/2026-04-08-001-refactor-web-workspace-boundary-plan.md` 标记为 `completed`
8. 仓库不再依赖“人记得边界”，而是由文档 + CI + 命令合同共同维持

## Self-Review

### Spec coverage

- 支持面定稿，有 Task 1
- guardrail，有 Task 2
- 工具链真相统一，有 Task 3
- 发布验证封口和 baseline 关闭，有 Task 4

### Placeholder scan

- 没有 `TBD`
- 没有“稍后处理错误”这类空话
- 所有关键命令都给了具体路径和预期

### Type consistency

- canonical / compatibility / repo contracts 三个术语全文保持一致
- 版本真相统一使用 Go 1.25+ / Bun 1.3.0+

## Execution Result

Task 1-4 已完成，`boundary guardrail`、`repo-contracts` 和 plan closeout 都已经落地。当前基线收口以 README / AGENTS 的支持面与版本真相、CI 的边界护栏和 `repo-contracts` canonical gate 为准，compatibility 和少量 deferred 项仍保留显式人工尾巴。
