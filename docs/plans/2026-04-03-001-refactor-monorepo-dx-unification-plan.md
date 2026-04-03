---
title: "refactor: Monorepo 开发体验统一"
type: refactor
status: completed
date: 2026-04-03
origin: docs/brainstorms/2026-04-03-monorepo-dx-unification-requirements.md
---

# refactor: Monorepo 开发体验统一

## Overview

统一 perfect-panel monorepo 的开发工具链：补齐根 Makefile 目标（build/format/typecheck/clean）、用根级 lefthook 替换分裂的 lefthook+husky 双系统、统一 commitlint 到 conventional commits、创建 server 配置示例文档。

## Problem Frame

当前仓库已建立基本 monorepo 结构（根 Makefile + Turbo workspace），但开发体验割裂：两套 git hooks 系统、Makefile 缺少常用目标、server 配置无文档。这导致开发者需手动记忆工具链差异，也为后续 CI/CD 留下不一致的基础。(see origin: docs/brainstorms/2026-04-03-monorepo-dx-unification-requirements.md)

## Requirements Trace

- R1. `make build` 同时构建 server 二进制和 web 应用
- R2. `make format` 同时执行 server go fmt/goimports 和 web biome format
- R3. `make typecheck` 执行 web 的 turbo check-types
- R4. `make clean` 清理构建产物（bin/、.next/、.turbo/）
- R5. 所有新增目标注册为 .PHONY
- R6. 根目录 `lefthook.yml` 统管 server 和 web 的 pre-commit 与 commit-msg
- R7. Server pre-commit 保留 go fmt/goimports/golangci-lint/go vet/go test，并行执行
- R8. Web pre-commit 对暂存文件执行 biome check --write（复现 lint-staged 行为）
- R9. commit-msg 统一使用 commitlint + conventional commits 规范
- R10. 移除 husky/lint-staged 依赖、配置、prepare.sh 脚本
- R11. 移除 server/lefthook.yml
- R12. bootstrap 目标包含 lefthook install
- R13. 创建 server/etc/ppanel.yaml.example
- R14. 配置示例覆盖 Config 全部 23 个顶层字段

## Scope Boundaries

- 不涉及 CI/CD pipeline 建设
- 不涉及 Docker Compose 统一
- 不涉及自动依赖更新（renovate/dependabot）
- 不涉及 web 测试框架引入
- 不改变运行时行为或部署拓扑

## Context & Research

### Relevant Code and Patterns

- `Makefile` — 现有目标使用 `cd <dir> && <cmd>` 模式，`SHELL := /bin/sh`，`.PHONY` 声明
- `server/lefthook.yml` — 5 个 parallel pre-commit 命令 + commitlint commit-msg
- `web/package.json` — husky/lint-staged/commitlint 配置均内联在此文件
- `web/packages/commitlint-config/base.js` — 当前扩展 gitmoji 规则
- `server/.commitlintrc.json` — 当前扩展 @commitlint/config-conventional
- `web/scripts/prepare.sh` — husky 初始化 + lobe 工具全局安装 + 交互式 lobe-commit（有 .git 检测 bug）
- `server/internal/config/config.go` — Config 结构体 23 个顶层字段，File 结构体 8 个字段

### External References

- [Biome Git Hooks 官方文档](https://biomejs.dev/recipes/git-hooks/) — 推荐 lefthook + `{staged_files}` + `stage_fixed: true`
- [lefthook monorepo 支持](https://github.com/evilmartians/lefthook/discussions/852) — `root` 选项指定子目录
- [lefthook commitlint 示例](https://lefthook.dev/examples/commitlint/) — commit-msg 使用 `{1}` 占位符

## Key Technical Decisions

- **Conventional Commits 统一**: 移除 gitmoji 配置，统一使用 `@commitlint/config-conventional`。与 CLAUDE.md 约定的 `<type>: <description>` 格式一致
- **lefthook staged_files + stage_fixed**: Web pre-commit 使用 `{staged_files}` 模板变量传递暂存文件列表，`stage_fixed: true` 自动 re-stage 被 biome 修改的文件。精确复现 lint-staged 行为
- **lefthook root 选项**: Server 命令使用 `root: "server/"`，Web 命令使用 `root: "web/"`，确保命令在正确子目录执行
- **commitlint 从 web workspace 执行**: 根级 lefthook 通过 `cd web && bunx commitlint --edit ../{1}` 调用（`{1}` ���开为 `.git/COMMIT_EDITMSG`，cd 到 web/ 后需 `../` 前缀回到仓库根），因 commitlint 及其配置安装在 web/node_modules
- **make clean 仅清理构建产物**: 清理 `server/bin/`、`web/**/. next/`、`web/**/.turbo/`，不删除 node_modules（恢复代价高）
- **make build server 策略**: 直接 `go build -o bin/ppanel-server .`（当前平台），不调用交叉编译 Makefile
- **删除 prepare.sh**: 整个脚本移除（含 .git 检测 bug 和交互式 lobe-commit），lobe 工具改为手动安装
- **ppanel.yaml 加入 .gitignore**: 防止开发者意外提交含真实密钥的配置，仅 .example 纳入版本控制

## Open Questions

### Resolved During Planning

- **lefthook 能否复现 lint-staged 行为？** 能。`{staged_files}` + `glob` 过滤 + `stage_fixed: true` 完全等效。Biome 官方推荐此方案
- **commitlint 从根目录如何执行？** `cd web && bunx commitlint --edit ../{1}`。commitlint 依赖安装在 web workspace，`{1}` 展开为 `.git/COMMIT_EDITMSG`（相对于仓库根），cd 到 web/ 后需 `../` 前缀回溯
- **prepare.sh 剩余内容如何处理？** 整个脚本删除。lobe-i18n/lobe-commit 为 i18n 翻译辅助工具，非核心开发依赖
- **commitlint 配置用哪套？** Conventional Commits。与 CLAUDE.md 约定一致，移除 gitmoji
- **make clean 是否删除 node_modules？** 不删。仅清理 .next/.turbo/bin/ 等构建产物

### Deferred to Implementation

- lefthook 的 `glob` 路径在 `root` 选项下的精确匹配行为，可能需要在实际安装后微调 glob 模式
- `web/packages/commitlint-config` 切换到 conventional 后是否需要调整 `semantic-release-config-gitmoji` 的 release 配置

## Implementation Units

- [ ] **Unit 1: 统一 commitlint 到 Conventional Commits**

**Goal:** 消除 commitlint 配置冲突，统一到 conventional commits 规范

**Requirements:** R9

**Dependencies:** None

**Files:**
- Modify: `web/packages/commitlint-config/base.js`
- Modify: `web/packages/commitlint-config/package.json`
- Delete: `server/.commitlintrc.json`

**Approach:**
- 将 `web/packages/commitlint-config/base.js` 的 `extends` 从 `["gitmoji"]` 改为 `["@commitlint/config-conventional"]`
- 在 `web/packages/commitlint-config/package.json` 中替换 `commitlint-config-gitmoji` 依赖为 `@commitlint/config-conventional`
- 删除 `server/.commitlintrc.json`，根级 lefthook 将通过 web workspace 调用统一的配置
- 保留 `base.js` 中的 `footer-leading-blank` 和 `header-max-length` 自定义规则（已禁用，无害）

**Patterns to follow:**
- `web/packages/commitlint-config/` 现有包结构

**Test scenarios:**
- Happy path: `feat: add new feature` 格式的提交信息通过 commitlint 校验
- Happy path: `fix(server): resolve crash` 带 scope 的格式通过校验
- Error path: `:sparkles: add feature` gitmoji 格式被拒绝
- Error path: 空提交信息被拒绝
- Edge case: `feat!: breaking change` BREAKING CHANGE 标记通过校验

**Verification:**
- `cd web && bunx commitlint --edit` 对 conventional 格式消息返回成功
- `server/.commitlintrc.json` 不再存在

---

- [ ] **Unit 2: 创建根级 lefthook.yml**

**Goal:** 统一 git hooks 管理，用一个根级配置替代分裂的两套系统

**Requirements:** R6, R7, R8, R9, R11

**Dependencies:** Unit 1（commitlint 配置就绪）

**Files:**
- Create: `lefthook.yml`
- Delete: `server/lefthook.yml`

**Approach:**
- 创建根目录 `lefthook.yml`，使用 lefthook 的 `root` 选项区分 server 和 web 命令
- Server pre-commit: 5 个命令（go-fmt、go-imports、go-lint、go-vet、go-test）设为 `parallel: true`，使用 `root: "server/"` 和 `glob: "*.go"` 确保仅在有 Go 文件变更时触发
- Web pre-commit: biome check 使用 `root: "web/"`、`glob: "*.{json,md,js,jsx,cjs,mjs,ts,tsx}"`、`{staged_files}` 模板、`stage_fixed: true`
- commit-msg: `cd web && bunx commitlint --edit ../{1}`（`{1}` 展开为 `.git/COMMIT_EDITMSG`，需 `../` 前缀从 web/ 解析回仓库根）
- 删除 `server/lefthook.yml`

**Patterns to follow:**
- `server/lefthook.yml` 现有命令结构（保留相同的命令名和参数）
- [Biome 官方 lefthook 配置](https://biomejs.dev/recipes/git-hooks/)

**Test scenarios:**
- Happy path: 修改 server/*.go 文件后 commit，触发 go fmt/goimports/golangci-lint/go vet/go test
- Happy path: 修改 web/*.ts 文件后 commit，触发 biome check --write 且修改后的文件自动 re-stage
- Happy path: commit-msg hook 使用 conventional commits 规则校验
- Edge case: 同时修改 server 和 web 文件，两组 hooks 均触发
- Edge case: 仅修改 web/*.md 文件，server hooks 不触发
- Error path: biome check 发现 lint 错误且无法自动修复时，commit 被阻止
- Integration: `{1}` 占位符正确展开为 .git/COMMIT_EDITMSG 路径

**Verification:**
- `lefthook run pre-commit` 从根目录执行成功
- `lefthook run commit-msg` 从根目录执行成功
- `server/lefthook.yml` 不再存在

---

- [ ] **Unit 3: 移除 husky/lint-staged/prepare.sh**

**Goal:** 清理被 lefthook 替代的旧 hooks 系统，消除死依赖

**Requirements:** R10

**Dependencies:** Unit 2（lefthook 已就位作为替代）

**Files:**
- Modify: `web/package.json`
- Modify: `web/.gitignore`
- Delete: `web/scripts/prepare.sh`

**Approach:**
- 从 `web/package.json` 中移除:
  - `devDependencies` 中的 `husky` 和 `lint-staged`
  - `lint-staged` 配置块（`"lint-staged": { "*.{json,md,...}": [...] }`）
  - `scripts.prepare`（指向 `./scripts/prepare.sh`）
- 删除 `web/scripts/prepare.sh` 整个文件
- 从 `web/.gitignore` 中移除 `.husky` 条目
- 运行 `cd web && bun install` 更新 lockfile

**Patterns to follow:**
- `web/package.json` 现有结构

**Test scenarios:**
- Happy path: `bun install` 在 web/ 中正常完成，无 prepare 脚本触发
- Happy path: `git grep -r "husky" web/` 返回空
- Edge case: 验证 web/package.json 中不再有 lint-staged 配置块
- Integration: lefthook（Unit 2）接管了 web pre-commit 和 commit-msg 功能

**Verification:**
- `web/scripts/prepare.sh` 不存在
- `web/package.json` 不含 husky、lint-staged 引用
- `bun install` 从 web/ 运行无报错

---

- [ ] **Unit 4: 补齐根 Makefile**

**Goal:** 提供统一的 make 命令入口，覆盖 build/format/typecheck/clean/bootstrap(hooks)

**Requirements:** R1, R2, R3, R4, R5, R12

**Dependencies:** Unit 2（bootstrap 需要 lefthook install）

**Files:**
- Modify: `Makefile`

**Approach:**
- 新增目标及子目标：
  - `build: server-build web-build`
    - `server-build`: `cd server && go build -o bin/ppanel-server .`
    - `web-build`: `cd web && bun run build`
  - `format: server-format web-format`
    - `server-format`: `cd server && go fmt ./... && goimports -w .`
    - `web-format`: `cd web && bun run format`
  - `typecheck: web-typecheck`
    - `web-typecheck`: `cd web && bun run typecheck`
  - `clean: server-clean web-clean`
    - `server-clean`: `rm -rf server/bin/`
    - `web-clean`: `cd web && rm -rf apps/*/.next apps/*/.turbo .turbo`
- 修改 `bootstrap` 目标：追加条件式 lefthook install（`command -v lefthook >/dev/null 2>&1 && lefthook install || echo "Warning: lefthook not found"`）
- 所有新增目标加入 `.PHONY` 声明
- 保持与现有目标一致的 `cd <dir> && <cmd>` 风格

**Patterns to follow:**
- `Makefile` 现有 `server-lint`/`web-lint` 的命名和结构

**Test scenarios:**
- Happy path: `make build` 产出 `server/bin/ppanel-server` 和 `web/apps/*/.next/`
- Happy path: `make format` 格式化 server Go 文件和 web 前端文件
- Happy path: `make typecheck` 执行 TypeScript 类型检查
- Happy path: `make clean` 清理 server/bin/ 和 web 的 .next/.turbo 目录
- Happy path: `make bootstrap` 安装依赖并执行 lefthook install
- Error path: `make build` 在 Go 编译错误时返回非零退出码
- Edge case: `make clean` 在无构建产物时幂等执行（不报错）

**Verification:**
- 从根目录执行 `make bootstrap && make lint && make format && make typecheck && make build && make clean` 全部成功
- `make bootstrap` 后 `.git/hooks/` 目录包含 lefthook 管理的 hook 文件

---

- [ ] **Unit 5: Server 配置示例文档**

**Goal:** 降低新贡献者入门门槛，提供完整的配置参考

**Requirements:** R13, R14

**Dependencies:** None

**Files:**
- Create: `server/etc/ppanel.yaml.example`
- Modify: `server/.gitignore`（或 `server/etc/.gitignore`）

**Approach:**
- 从 `server/internal/config/config.go` 的 Config 结构体提取全部 23 个顶层字段
- 为每个字段提供 YAML 键名、类型说明（注释）和合理的默认值或占位值
- 嵌套结构体（如 JwtAuth、MySQL、Redis 等）展开为 YAML 层级，包含子字段
- 将 `etc/ppanel.yaml` 加入 `server/.gitignore`
- 从 git 中取消跟踪现有空文件：需执行 `git rm --cached server/etc/ppanel.yaml`

**Patterns to follow:**
- `web/apps/admin/.env.template` 和 `web/apps/user/.env.template` 的文档风格（注释说明每个变量）
- `server/internal/config/config.go` 中 Config 和 File 结构体的字段定义

**Test scenarios:**
- Happy path: 复制 ppanel.yaml.example 为 ppanel.yaml，填入真实 DB/Redis 地址后 `go run . run --config etc/ppanel.yaml` 正常启动
- Edge case: ppanel.yaml.example 中所有必填字段（MySQL.Addr、Redis.Host、JwtAuth.AccessSecret）有占位值和注释标注
- Error path: 缺少必填字段时 server 启动进入初始化流程（而非崩溃）

**Verification:**
- `server/etc/ppanel.yaml.example` 存在且包含全部 23 个顶层字段
- `server/etc/ppanel.yaml` 已从 git 跟踪中移除
- `git status` 不再显示 ppanel.yaml 为已跟踪文件

## System-Wide Impact

- **Interaction graph:** lefthook 替代 husky 后，`bun install` 不再触发 prepare 脚本（hooks 由 `make bootstrap` 中的 `lefthook install` 管理）
- **Error propagation:** lefthook pre-commit 失败会阻止 commit，行为与当前 husky 一致
- **State lifecycle risks:** `git rm --cached server/etc/ppanel.yaml` 会在队友 pull 后从其工作目录删除该文件——需在 commit message 中明确说明此变更
- **API surface parity:** commitlint 从 gitmoji 切换到 conventional 会拒绝现有 gitmoji 格式的提交，影响所有开发者
- **Unchanged invariants:** server 和 web 的运行时行为、构建产物格式、部署拓扑不受影响

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| lefthook glob 路径在 root 选项下匹配不精确 | 安装后手动测试，微调 glob 模式 |
| 队友 pull 后 ppanel.yaml 被删除 | commit message 中明确说明，README 更新 |
| commitlint 切换影响现有 gitmoji 提交者 | 在 PR 描述中说明规范变更 |
| semantic-release-config-gitmoji 可能与 conventional commits 不兼容 | 此项属于 release 流程，超出当前 scope，后续单独处理 |
| lefthook CLI 未安装时 `make bootstrap` 中 `lefthook install` 失败 | bootstrap 可使用 `command -v lefthook && lefthook install` 做条件安装，或在 README 中标明前置依赖 |

## Sources & References

- **Origin document:** [docs/brainstorms/2026-04-03-monorepo-dx-unification-requirements.md](docs/brainstorms/2026-04-03-monorepo-dx-unification-requirements.md)
- Related code: `Makefile`, `server/lefthook.yml`, `web/package.json`, `web/scripts/prepare.sh`, `server/internal/config/config.go`
- External docs: [Biome Git Hooks](https://biomejs.dev/recipes/git-hooks/), [lefthook monorepo](https://github.com/evilmartians/lefthook/discussions/852)
