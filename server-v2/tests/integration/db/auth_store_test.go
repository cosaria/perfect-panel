package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
	authstore "github.com/perfect-panel/server-v2/internal/domains/auth/store"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func TestAuthSQLStoreConsumesPasswordResetOnceAtomically(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	if err := ppdb.Migrate(context.Background(), db); err != nil {
		t.Fatalf("执行 migrate 失败: %v", err)
	}

	const (
		userID     = "00000000-0000-4000-8000-000000000101"
		identityID = "00000000-0000-4000-8000-000000000102"
		tokenID    = "00000000-0000-4000-8000-000000000103"
	)

	if _, err := db.Exec(
		`INSERT INTO users(id, status) VALUES ($1, 'active')`,
		userID,
	); err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO user_identities(id, user_id, provider, identifier, secret_hash, status)
		VALUES ($1, $2, 'email', 'user@example.com', 'hashed:old-password', 'active')`,
		identityID,
		userID,
	); err != nil {
		t.Fatalf("插入身份失败: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO verification_tokens(id, user_id, purpose, token_hash, destination, expires_at, status)
		VALUES ($1, $2, 'password_reset', 'token:reset-plain', 'user@example.com', $3, 'pending')`,
		tokenID,
		userID,
		time.Date(2026, 4, 10, 18, 0, 0, 0, time.UTC),
	); err != nil {
		t.Fatalf("插入验证码失败: %v", err)
	}

	store := authstore.NewSQLStore(db)
	usedAt := time.Date(2026, 4, 10, 10, 30, 0, 0, time.UTC)
	events := []authdomain.AuthEvent{
		{
			ID:        "00000000-0000-4000-8000-000000000104",
			Type:      authdomain.AuthEventTypePasswordResetCompleted,
			CreatedAt: usedAt,
			Payload: map[string]any{
				"tokenHash": "token:reset-plain",
			},
		},
	}

	record, err := store.ConsumePasswordReset(context.Background(), "token:reset-plain", "hashed:new-password", usedAt, events)
	if err != nil {
		t.Fatalf("第一次消费密码重置 token 失败: %v", err)
	}
	if record.Token.ID != tokenID {
		t.Fatalf("返回 token_id 不匹配: want=%s got=%s", tokenID, record.Token.ID)
	}
	if record.Identity.SecretHash != "hashed:new-password" {
		t.Fatalf("返回 identity.secret_hash 不匹配: got=%s", record.Identity.SecretHash)
	}

	var storedUsedAt sql.NullTime
	var tokenStatus string
	if err := db.QueryRow(
		`SELECT used_at, status FROM verification_tokens WHERE id = $1`,
		tokenID,
	).Scan(&storedUsedAt, &tokenStatus); err != nil {
		t.Fatalf("查询验证码状态失败: %v", err)
	}
	if !storedUsedAt.Valid || !storedUsedAt.Time.Equal(usedAt) {
		t.Fatalf("verification_tokens.used_at 未正确落库: %+v", storedUsedAt)
	}
	if tokenStatus != authdomain.VerificationTokenStatusConsumed {
		t.Fatalf("verification_tokens.status 不匹配: got=%s", tokenStatus)
	}

	var secretHash string
	if err := db.QueryRow(
		`SELECT secret_hash FROM user_identities WHERE id = $1`,
		identityID,
	).Scan(&secretHash); err != nil {
		t.Fatalf("查询身份密码哈希失败: %v", err)
	}
	if secretHash != "hashed:new-password" {
		t.Fatalf("身份密码哈希未更新: got=%s", secretHash)
	}

	var eventCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM auth_events WHERE event_type = $1`, authdomain.AuthEventTypePasswordResetCompleted).Scan(&eventCount); err != nil {
		t.Fatalf("统计密码重置事件失败: %v", err)
	}
	if eventCount != 1 {
		t.Fatalf("应写入 1 条密码重置完成事件，实际 %d", eventCount)
	}

	_, err = store.ConsumePasswordReset(context.Background(), "token:reset-plain", "hashed:another-password", usedAt.Add(time.Minute), events)
	if !errors.Is(err, authdomain.ErrVerificationTokenConsumed) {
		t.Fatalf("重复消费同一 token 应返回 ErrVerificationTokenConsumed，实际 %v", err)
	}

	if err := db.QueryRow(
		`SELECT secret_hash FROM user_identities WHERE id = $1`,
		identityID,
	).Scan(&secretHash); err != nil {
		t.Fatalf("重复消费后查询身份密码哈希失败: %v", err)
	}
	if secretHash != "hashed:new-password" {
		t.Fatalf("重复消费不应覆盖第一次写入的密码哈希，got=%s", secretHash)
	}
}
