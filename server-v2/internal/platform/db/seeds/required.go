package seeds

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/perfect-panel/server-v2/internal/domains/access"
	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

// ApplyRequired 写入系统启动所需的最小种子数据。
func ApplyRequired(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	return ppdb.WithTx(ctx, database, func(tx *sql.Tx) error {
		for _, role := range access.RequiredRoles() {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO roles(code, name) VALUES ($1, $2)
				ON CONFLICT(code) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()`,
				role.Code,
				role.Name,
			); err != nil {
				return fmt.Errorf("写入基础角色失败: %w", err)
			}
		}

		for _, permission := range access.RequiredPermissions() {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO permissions(code, name) VALUES ($1, $2)
				ON CONFLICT(code) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()`,
				permission.Code,
				permission.Name,
			); err != nil {
				return fmt.Errorf("写入权限失败: %w", err)
			}
		}

		for roleCode, permissionCodes := range access.RequiredRolePermissions() {
			for _, permissionCode := range permissionCodes {
				if _, err := tx.ExecContext(
					ctx,
					`INSERT INTO role_permissions(role_id, permission_id)
					SELECT roles.id, permissions.id
					FROM roles
					JOIN permissions ON permissions.code = $2
					WHERE roles.code = $1
					ON CONFLICT(role_id, permission_id) DO NOTHING`,
					roleCode,
					permissionCode,
				); err != nil {
					return fmt.Errorf("写入角色权限关联失败: %w", err)
				}
			}
		}

		for _, setting := range authdomain.RequiredSettings() {
			raw, err := json.Marshal(setting.Value)
			if err != nil {
				return fmt.Errorf("序列化系统配置失败 %s/%s: %w", setting.Scope, setting.Key, err)
			}
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO system_settings(scope, key, value_json) VALUES ($1, $2, $3::jsonb)
				ON CONFLICT(scope, key) DO UPDATE SET value_json = EXCLUDED.value_json, updated_at = NOW()`,
				setting.Scope,
				setting.Key,
				string(raw),
			); err != nil {
				return fmt.Errorf("写入系统配置失败: %w", err)
			}
		}
		return nil
	})
}
