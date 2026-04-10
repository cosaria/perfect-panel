package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envDBDSN = "PPANEL_DB_DSN"
)

var (
	// ErrMissingDSN 表示未提供数据库连接串。
	ErrMissingDSN = errors.New("未设置数据库连接串环境变量 PPANEL_DB_DSN")
)

// OpenFromEnv 从环境变量读取 DSN 并建立连接。
func OpenFromEnv() (*sql.DB, error) {
	return Open(strings.TrimSpace(os.Getenv(envDBDSN)))
}

// Open 建立 PostgreSQL 连接并执行一次 ping 校验。
func Open(dsn string) (*sql.DB, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return nil, ErrMissingDSN
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("数据库 ping 失败: %w", err)
	}
	return db, nil
}
