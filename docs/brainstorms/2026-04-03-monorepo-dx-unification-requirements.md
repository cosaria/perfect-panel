---
date: 2026-04-03
topic: monorepo-dx-unification
---

# Monorepo 开发体验统一

## Problem Frame

perfect-panel 已在 `feat/monorepo-baseline` 分支建立了基本的 monorepo 结构（根 Makefile + Turbo workspace），但开发体验仍然割裂：根 Makefile 缺少 build/format/typecheck/clean 等常用目标；server 和 web 分别使用 lefthook 和 husky 两套独立的 git hooks 系统；server 配置文件缺少示例文档，新贡献者入门门槛高。

这些缺口导致开发者需要手动记忆不同子项目的工具链差异，无法通过统一入口完成日常操作，也为后续 CI/CD 建设留下不一致的基础。

## Requirements

**Root Makefile 补齐**

- R1. 新增 `make build` 目标，同时构建 server 二进制和 web 应用（通过 turbo build）
- R2. 新增 `make format` 目标，同时执行 server 的 `go fmt` + `goimports` 和 web 的 `biome format --write`
- R3. 新增 `make typecheck` 目标，执行 web 的 `turbo check-types`（Go 编译器已隐式做类型检查，无需额外目标）
- R4. 新增 `make clean` 目标，清理 server 的 `bin/` 和 web 的 `node_modules`/`.next` 等构建产物
- R5. 所有新增目标须注册为 `.PHONY`，与现有目标保持一致风格

**Git Hooks 统一（lefthook）**

- R6. 在仓库根目录创建 `lefthook.yml`，统管 server 和 web 的 pre-commit 与 commit-msg hooks
- R7. Server pre-commit hooks 保留现有全部检查：`go fmt`、`goimports`、`golangci-lint`、`go vet`、`go test -v ./...`，并行执行
- R8. Web pre-commit hooks 复现当前 lint-staged 行为：对暂存文件执行 `biome check --write`
- R9. commit-msg hook 统一使用 commitlint 校验提交格式
- R10. 移除 web/package.json 中的 husky 和 lint-staged devDependencies 及 lint-staged 配置块；移除 `web/scripts/prepare.sh` 中的 husky 初始化逻辑；移除 web/.gitignore 中的 `.husky` 条目
- R11. 移除 `server/lefthook.yml`（配置已提升到根级）
- R12. 根 Makefile 的 `bootstrap` 目标须包含 `lefthook install`，确保 clone 后自动安装 hooks

**Server 配置文档**

- R13. 创建 `server/etc/ppanel.yaml.example`，列出所有可配置字段及合理的默认值/注释说明
- R14. 配置示例须覆盖 Config 结构体全部顶层字段：Model、Host、Port、Debug、TLS、JwtAuth、Logger、MySQL、Redis、Site、Node、Mobile、Email、Device、Verify、VerifyCode、Register、Subscribe、Invite、Telegram、Log、Currency、Administrator

## Success Criteria

- 从仓库根目录执行 `make bootstrap && make lint && make format && make typecheck && make build && make clean`，全部成功
- `git commit` 触发根级 lefthook，server 和 web 的检查均正常运行
- 仓库中不再存在 husky 相关配置和依赖
- 新贡献者复制 `ppanel.yaml.example` 为 `ppanel.yaml`，填入 DB/Redis 连接信息即可启动 server

## Scope Boundaries

- 不涉及 CI/CD pipeline 建设（留给后续迭代）
- 不涉及 Docker Compose 统一（留给后续迭代）
- 不涉及自动依赖更新（renovate/dependabot）
- 不涉及 web 测试框架引入
- 不改变现有 server/web 的运行时行为或部署拓扑

## Key Decisions

- **lefthook 而非 husky**: lefthook 原生支持 monorepo（root 配置可指定子目录 glob），Go 编写无 Node 依赖，与 server 侧已有实践一致
- **pre-commit 保留完整 go test**: 虽然慢，但用户选择确保每次提交都绿；后续可通过 CI 覆盖后再精简
- **typecheck 只覆盖 web**: Go 编译器已做类型检查，无需额外 Makefile 目标

## Dependencies / Assumptions

- lefthook CLI 已安装或可通过 `go install` / `brew install` 获取
- goimports 需通过 `go install golang.org/x/tools/cmd/goimports@latest` 安装
- `commitlint` 及 `@workspace/commitlint-config` 在 web workspace 中已存在，从根级 lefthook 可通过 `cd web && npx commitlint` 调用

## Outstanding Questions

### Deferred to Planning

- [Affects R8][Technical] lefthook 对 web 暂存文件的 `biome check --write` 如何精确复现 lint-staged 的行为（只检查暂存文件而非全量）？需研究 lefthook 的 `glob`/`staged_files` 配置
- [Affects R10][Needs research] `web/scripts/prepare.sh` 中除 husky 外还有 `lobe-commit` 和 `lobe-i18n` 的全局安装，移除 husky 逻辑后如何处理剩余内容？
## Next Steps

→ `/ce:plan` for structured implementation planning
