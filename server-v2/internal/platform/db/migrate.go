package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

var managedTables = []string{
	"users",
	"roles",
	"system_settings",
	"outbox_events",
	"schema_revisions",
}

// Migrate 执行当前 baseline migration，并记录目标 schema version。
func Migrate(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	revisionExists, err := schemaRevisionExists(ctx, database, TargetSchemaVersion)
	if err != nil {
		return err
	}
	if revisionExists {
		return nil
	}

	existingTables, err := findExistingManagedTables(ctx, database)
	if err != nil {
		return err
	}
	if len(existingTables) > 0 {
		return fmt.Errorf("检测到 schema 已漂移或状态异常，拒绝执行 baseline migration: %s", strings.Join(existingTables, ","))
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
			`INSERT INTO schema_revisions(version) VALUES ($1)`,
			TargetSchemaVersion,
		); err != nil {
			return fmt.Errorf("写入 schema_revisions 失败: %w", err)
		}
		return nil
	})
}

func schemaRevisionExists(ctx context.Context, database *sql.DB, version string) (bool, error) {
	var revisionTableExists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT to_regclass(current_schema() || '.schema_revisions') IS NOT NULL`,
	).Scan(&revisionTableExists); err != nil {
		return false, fmt.Errorf("检查 schema_revisions 是否存在失败: %w", err)
	}
	if !revisionTableExists {
		return false, nil
	}

	var count int
	if err := database.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM schema_revisions WHERE version = $1`,
		version,
	).Scan(&count); err != nil {
		return false, fmt.Errorf("检查目标 schema version 失败: %w", err)
	}
	return count > 0, nil
}

func findExistingManagedTables(ctx context.Context, database *sql.DB) ([]string, error) {
	inClauseParts := make([]string, 0, len(managedTables))
	for _, table := range managedTables {
		inClauseParts = append(inClauseParts, `'`+strings.ReplaceAll(table, `'`, `''`)+`'`)
	}

	rows, err := database.QueryContext(
		ctx,
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = current_schema() AND table_name IN (`+strings.Join(inClauseParts, ",")+`)`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询受管表状态失败: %w", err)
	}
	defer rows.Close()

	existing := make([]string, 0, len(managedTables))
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("扫描受管表名称失败: %w", err)
		}
		existing = append(existing, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历受管表状态失败: %w", err)
	}
	slices.Sort(existing)
	return existing, nil
}
