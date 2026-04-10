package seeds

import (
	"context"
	"database/sql"
	"fmt"

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
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO roles(code, name) VALUES
				('admin', '管理员'),
				('user', '普通用户')
			ON CONFLICT(code) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()`,
		); err != nil {
			return fmt.Errorf("写入基础角色失败: %w", err)
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO system_settings(key, value) VALUES
				('site_name', 'Perfect Panel'),
				('app_name', 'Perfect Panel')
			ON CONFLICT(key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()`,
		); err != nil {
			return fmt.Errorf("写入系统配置失败: %w", err)
		}
		return nil
	})
}
