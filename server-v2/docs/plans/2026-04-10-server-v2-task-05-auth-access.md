# Task 05: 实现 `auth` 与 `access` 主链

**状态：** `未开始`  
**前置任务：** Task 04  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 5`

## 目标

打通 `public / user / admin` 调用面所必需的主体、身份、会话、验证码、密码重置和 RBAC 主链，并把 `required seed` 扩成可运行的权限最小集。

## 执行边界

- 允许触达：
  - `server-v2/internal/domains/auth/**`
  - `server-v2/internal/domains/access/**`
  - `server-v2/internal/platform/http/middleware/{session_auth.go,access_guard.go}`
  - `server-v2/internal/app/routing/{public.go,user.go,admin.go}`
  - `server-v2/internal/platform/db/migrations/0002_auth_access.sql`
  - `server-v2/internal/platform/db/seeds/required.go`
  - `server-v2/cmd/server/seed_required.go`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/public/{sessions,verification_tokens,password_reset_requests,password_resets}.yaml`
  - `server-v2/openapi/paths/user/me_sessions.yaml`
  - `server-v2/tests/domains/auth/usecase/*`
  - `server-v2/tests/contract/auth_api_contract_test.go`

## 关键产物

- `User / Identity / Session / VerificationToken / AuthEvent`
- `Role / Permission / UserRole / RolePermission`
- `sign_in / sign_out / issue_verification / request_password_reset / reset_password`
- `session_auth.go`
- `access_guard.go`
- `seed_permissions.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/domains/auth/usecase ./tests/contract -run 'TestSignInCreatesSession|TestResetPasswordConsumesVerificationTokenOnce|TestPublicSessionsContract|TestPasswordResetContract' -count=1
```

## 放行标准

- `verification_tokens` 只存哈希，不存明文。
- 密码重置和验证码消费必须是原子单次消费。
- `admin` 路由必须是 `RequireSession` 后再经过 `RequirePermissions`。
- `seed-required` 已补齐角色、权限、角色权限和认证默认项。

## 默认提交点

```bash
git add server-v2/internal/domains/auth server-v2/internal/domains/access server-v2/internal/platform/http/middleware server-v2/internal/app/routing server-v2/internal/platform/db/migrations/0002_auth_access.sql server-v2/internal/platform/db/seeds/required.go server-v2/cmd/server/seed_required.go server-v2/openapi/openapi.yaml server-v2/openapi/paths/public/sessions.yaml server-v2/openapi/paths/public/verification_tokens.yaml server-v2/openapi/paths/public/password_reset_requests.yaml server-v2/openapi/paths/public/password_resets.yaml server-v2/openapi/paths/user/me_sessions.yaml server-v2/tests/domains/auth/usecase server-v2/tests/contract/auth_api_contract_test.go
git commit -m "feat(server-v2): implement auth and access foundation"
```

## 完成后进入

- [2026-04-10-server-v2-task-06-system-catalog.md](./2026-04-10-server-v2-task-06-system-catalog.md)
