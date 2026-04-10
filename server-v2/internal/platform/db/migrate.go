package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
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
	"user_identities",
	"user_sessions",
	"verification_tokens",
	"permissions",
	"role_permissions",
	"user_roles",
	"auth_events",
}

var schemaMigrations = []string{
	"0001_baseline.sql",
	"0002_auth_access.sql",
}

// Migrate 执行 revision 主链，并记录目标 schema version。
func Migrate(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	revisionTableExists, currentVersion, err := currentSchemaVersionState(ctx, database)
	if err != nil {
		return err
	}
	if revisionTableExists {
		if currentVersion == "" {
			return fmt.Errorf("schema_revisions 状态异常: 当前版本为空")
		}
		currentIndex := revisionIndex(currentVersion)
		if currentIndex < 0 {
			return fmt.Errorf("当前 schema version 不在受管 revision 主链中: got=%s", currentVersion)
		}
		if currentVersion == TargetSchemaVersion {
			return EnsureSchemaVersion(ctx, database)
		}
		return applyMigrationsFromIndex(ctx, database, currentIndex+1)
	}

	existingTables, err := findExistingManagedTables(ctx, database)
	if err != nil {
		return err
	}
	if len(existingTables) > 0 {
		return fmt.Errorf("检测到 schema 已漂移或状态异常，拒绝执行 baseline migration: %s", strings.Join(existingTables, ","))
	}

	return applyMigrationsFromIndex(ctx, database, 0)
}

func currentSchemaVersionState(ctx context.Context, database *sql.DB) (bool, string, error) {
	var revisionTableExists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT to_regclass(current_schema() || '.schema_revisions') IS NOT NULL`,
	).Scan(&revisionTableExists); err != nil {
		return false, "", fmt.Errorf("检查 schema_revisions 是否存在失败: %w", err)
	}
	if !revisionTableExists {
		return false, "", nil
	}

	var version string
	if err := database.QueryRowContext(
		ctx,
		`SELECT version FROM schema_revisions ORDER BY applied_at DESC, id DESC LIMIT 1`,
	).Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, "", nil
		}
		return false, "", fmt.Errorf("读取当前 schema version 失败: %w", err)
	}
	return true, version, nil
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

func applyMigrationsFromIndex(ctx context.Context, database *sql.DB, start int) error {
	for idx := start; idx < len(schemaMigrations); idx++ {
		version := strings.TrimSuffix(filepath.Base(schemaMigrations[idx]), filepath.Ext(schemaMigrations[idx]))
		script, err := migrationFiles.ReadFile("migrations/" + schemaMigrations[idx])
		if err != nil {
			return fmt.Errorf("读取 migration %s 失败: %w", version, err)
		}
		if err := WithTx(ctx, database, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, string(script)); err != nil {
				return fmt.Errorf("执行 migration %s 失败: %w", version, err)
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO schema_revisions(version) VALUES ($1)`, version); err != nil {
				return fmt.Errorf("写入 schema_revisions %s 失败: %w", version, err)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return EnsureSchemaVersion(ctx, database)
}

func revisionIndex(version string) int {
	for idx, name := range schemaMigrations {
		current := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
		if current == version {
			return idx
		}
	}
	return -1
}
