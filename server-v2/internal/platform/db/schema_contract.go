package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ValidateSchemaContract 校验 baseline 最小承重契约。
func ValidateSchemaContract(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	for _, table := range managedTables {
		ok, err := tableExists(ctx, database, table)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("schema 契约缺失: 表 %s 不存在", table)
		}
	}

	if err := assertUUIDColumn(ctx, database, "users", "id"); err != nil {
		return err
	}
	if err := assertUUIDColumn(ctx, database, "roles", "id"); err != nil {
		return err
	}
	if err := assertUUIDDefault(ctx, database, "users", "id"); err != nil {
		return err
	}
	if err := assertUUIDDefault(ctx, database, "roles", "id"); err != nil {
		return err
	}
	if err := assertUUIDDefault(ctx, database, "system_settings", "id"); err != nil {
		return err
	}
	if err := assertColumnType(ctx, database, "users", "status", "text"); err != nil {
		return err
	}

	if err := assertColumnType(ctx, database, "system_settings", "scope", "text"); err != nil {
		return err
	}
	if err := assertColumnType(ctx, database, "system_settings", "key", "text"); err != nil {
		return err
	}
	if err := assertColumnType(ctx, database, "system_settings", "value_json", "jsonb"); err != nil {
		return err
	}

	if err := assertColumnExists(ctx, database, "outbox_events", "aggregate_type"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "outbox_events", "aggregate_id"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "user_identities", "provider"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "user_identities", "identifier"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "user_sessions", "token_hash"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "verification_tokens", "token_hash"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "verification_tokens", "used_at"); err != nil {
		return err
	}
	if err := assertColumnExists(ctx, database, "auth_events", "event_type"); err != nil {
		return err
	}

	ok, err := hasSystemSettingsScopeKeyUniqueIndex(ctx, database)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("schema 契约缺失: system_settings 需要 (scope, key) 唯一索引/约束")
	}
	if err := assertUniqueIndexContains(ctx, database, "verification_tokens", "(token_hash)"); err != nil {
		return err
	}
	if err := assertUniqueIndexContains(ctx, database, "user_identities", "(provider, identifier)"); err != nil {
		return err
	}
	return nil
}

func tableExists(ctx context.Context, database *sql.DB, table string) (bool, error) {
	var exists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT to_regclass(current_schema() || '.' || $1) IS NOT NULL`,
		table,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("检查表 %s 是否存在失败: %w", table, err)
	}
	return exists, nil
}

func assertUUIDColumn(ctx context.Context, database *sql.DB, table string, column string) error {
	return assertColumnType(ctx, database, table, column, "uuid")
}

func assertColumnExists(ctx context.Context, database *sql.DB, table string, column string) error {
	var exists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2
		)`,
		table,
		column,
	).Scan(&exists); err != nil {
		return fmt.Errorf("检查列 %s.%s 是否存在失败: %w", table, column, err)
	}
	if !exists {
		return fmt.Errorf("schema 契约缺失: 列 %s.%s 不存在", table, column)
	}
	return nil
}

func assertColumnType(ctx context.Context, database *sql.DB, table string, column string, wantType string) error {
	var udtName string
	err := database.QueryRowContext(
		ctx,
		`SELECT udt_name
		FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2`,
		table,
		column,
	).Scan(&udtName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("schema 契约缺失: 列 %s.%s 不存在", table, column)
		}
		return fmt.Errorf("读取列类型 %s.%s 失败: %w", table, column, err)
	}
	if strings.ToLower(udtName) != strings.ToLower(wantType) {
		return fmt.Errorf("schema 契约不匹配: 列 %s.%s 期望类型 %s，实际 %s", table, column, wantType, udtName)
	}
	return nil
}

func assertUUIDDefault(ctx context.Context, database *sql.DB, table string, column string) error {
	var columnDefault sql.NullString
	err := database.QueryRowContext(
		ctx,
		`SELECT column_default
		FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2`,
		table,
		column,
	).Scan(&columnDefault)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("schema 契约缺失: 列 %s.%s 不存在", table, column)
		}
		return fmt.Errorf("读取列默认值 %s.%s 失败: %w", table, column, err)
	}
	if !columnDefault.Valid {
		return fmt.Errorf("schema 契约不匹配: 列 %s.%s 缺少 gen_random_uuid() 默认值", table, column)
	}
	if !strings.Contains(strings.ToLower(columnDefault.String), "gen_random_uuid()") {
		return fmt.Errorf("schema 契约不匹配: 列 %s.%s 默认值应包含 gen_random_uuid()，实际 %s", table, column, columnDefault.String)
	}
	return nil
}

func hasSystemSettingsScopeKeyUniqueIndex(ctx context.Context, database *sql.DB) (bool, error) {
	var exists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = current_schema()
				AND tablename = 'system_settings'
				AND position('create unique index' IN lower(indexdef)) > 0
				AND position('(scope, key)' IN lower(replace(indexdef, '"', ''))) > 0
		)`,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("检查 system_settings 唯一索引失败: %w", err)
	}
	return exists, nil
}

func assertUniqueIndexContains(ctx context.Context, database *sql.DB, table string, snippet string) error {
	var exists bool
	if err := database.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = current_schema()
				AND tablename = $1
				AND position('create unique index' IN lower(indexdef)) > 0
				AND position($2 IN lower(replace(indexdef, '"', ''))) > 0
		)`,
		table,
		strings.ToLower(snippet),
	).Scan(&exists); err != nil {
		return fmt.Errorf("检查 %s 唯一索引失败: %w", table, err)
	}
	if !exists {
		return fmt.Errorf("schema 契约缺失: %s 需要唯一索引 %s", table, snippet)
	}
	return nil
}
