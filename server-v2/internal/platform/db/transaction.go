package db

import (
	"context"
	"database/sql"
	"fmt"
)

// WithTx 用最小模板包装事务生命周期。
func WithTx(ctx context.Context, database *sql.DB, fn func(tx *sql.Tx) error) error {
	if database == nil {
		return fmt.Errorf("数据库实例不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}
