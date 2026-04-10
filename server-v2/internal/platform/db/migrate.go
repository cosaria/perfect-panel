package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate 执行当前 baseline migration，并记录目标 schema version。
func Migrate(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	script, err := migrationFiles.ReadFile("migrations/0001_baseline.sql")
	if err != nil {
		return fmt.Errorf("读取 baseline migration 失败: %w", err)
	}

	return WithTx(ctx, database, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, string(script)); err != nil {
			return fmt.Errorf("执行 baseline migration 失败: %w", err)
		}
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO schema_revisions(version) VALUES ($1) ON CONFLICT(version) DO NOTHING`,
			TargetSchemaVersion,
		); err != nil {
			return fmt.Errorf("写入 schema_revisions 失败: %w", err)
		}
		return nil
	})
}
