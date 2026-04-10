package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	// TargetSchemaVersion 表示当前代码要求的目标 schema 版本。
	TargetSchemaVersion = "0001_baseline"
)

// CurrentSchemaVersion 返回当前 schema 的最近版本号。
func CurrentSchemaVersion(ctx context.Context, database *sql.DB) (string, error) {
	if database == nil {
		return "", fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var tableExists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT to_regclass(current_schema() || '.schema_revisions') IS NOT NULL`,
	).Scan(&tableExists); err != nil {
		return "", fmt.Errorf("查询 schema_revisions 是否存在失败: %w", err)
	}
	if !tableExists {
		return "", fmt.Errorf("schema_revisions 表不存在")
	}

	var version string
	err := database.QueryRowContext(
		ctx,
		`SELECT version FROM schema_revisions ORDER BY applied_at DESC, id DESC LIMIT 1`,
	).Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("schema_revisions 为空")
		}
		return "", fmt.Errorf("读取 schema version 失败: %w", err)
	}
	return version, nil
}

// EnsureSchemaVersion 校验当前 schema 是否匹配目标版本。
func EnsureSchemaVersion(ctx context.Context, database *sql.DB) error {
	current, err := CurrentSchemaVersion(ctx, database)
	if err != nil {
		return err
	}
	if current != TargetSchemaVersion {
		return fmt.Errorf("schema version 不匹配: want=%s got=%s", TargetSchemaVersion, current)
	}
	return ValidateSchemaContract(ctx, database)
}
