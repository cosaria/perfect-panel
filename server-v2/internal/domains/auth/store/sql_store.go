package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	authmodel "github.com/perfect-panel/server-v2/internal/domains/auth/model"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) FindIdentityByIdentifier(ctx context.Context, provider string, identifier string) (authmodel.Identity, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, provider, identifier, secret_hash, status, verified_at, created_at, updated_at
		FROM user_identities
		WHERE provider = $1 AND identifier = $2
	`, provider, identifier)
	return scanIdentity(row)
}

func (s *SQLStore) FindUserByID(ctx context.Context, userID string) (authmodel.User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, status, archived_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID)
	return scanUser(row)
}

func (s *SQLStore) CreateSession(ctx context.Context, session authmodel.Session, events []authmodel.AuthEvent) error {
	return ppdb.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_sessions(id, user_id, token_hash, last_seen_at, expires_at, revoked_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, session.ID, session.UserID, session.TokenHash, session.LastSeenAt, session.ExpiresAt, nullableTime(session.RevokedAt), session.CreatedAt, session.UpdatedAt); err != nil {
			return fmt.Errorf("创建会话失败: %w", err)
		}
		return insertAuthEventsTx(ctx, tx, events)
	})
}

func (s *SQLStore) SaveVerificationToken(ctx context.Context, token authmodel.VerificationToken, events []authmodel.AuthEvent) error {
	return ppdb.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO verification_tokens(id, user_id, purpose, token_hash, destination, expires_at, used_at, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, token.ID, nullString(token.UserID), token.Purpose, token.TokenHash, token.Destination, token.ExpiresAt, nullableTime(token.UsedAt), token.Status, token.CreatedAt, token.UpdatedAt); err != nil {
			return fmt.Errorf("保存验证码失败: %w", err)
		}
		return insertAuthEventsTx(ctx, tx, events)
	})
}

func (s *SQLStore) ConsumePasswordReset(ctx context.Context, tokenHash string, passwordHash string, usedAt time.Time, events []authmodel.AuthEvent) (authmodel.PasswordResetRecord, error) {
	var record authmodel.PasswordResetRecord
	err := ppdb.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		token, err := consumePasswordResetToken(ctx, tx, tokenHash, usedAt)
		if err != nil {
			return err
		}

		identityRow := tx.QueryRowContext(ctx, `
			UPDATE user_identities
			SET secret_hash = $1, updated_at = $2
			WHERE user_id = $3 AND provider = $4 AND identifier = $5
			RETURNING id, user_id, provider, identifier, secret_hash, status, verified_at, created_at, updated_at
		`, passwordHash, usedAt, token.UserID, authmodel.IdentityProviderEmail, token.Destination)
		identity, err := scanIdentity(identityRow)
		if err != nil {
			if errors.Is(err, authmodel.ErrIdentityNotFound) {
				return authmodel.ErrIdentityNotFound
			}
			return err
		}

		record = authmodel.PasswordResetRecord{
			Token:    token,
			Identity: identity,
		}
		enrichedEvents := append([]authmodel.AuthEvent(nil), events...)
		for idx := range enrichedEvents {
			if enrichedEvents[idx].UserID == "" {
				enrichedEvents[idx].UserID = identity.UserID
			}
			if enrichedEvents[idx].Payload == nil {
				enrichedEvents[idx].Payload = map[string]any{}
			}
			if _, ok := enrichedEvents[idx].Payload["tokenId"]; !ok {
				enrichedEvents[idx].Payload["tokenId"] = token.ID
			}
		}
		return insertAuthEventsTx(ctx, tx, enrichedEvents)
	})
	if err != nil {
		return authmodel.PasswordResetRecord{}, err
	}
	return record, nil
}

func (s *SQLStore) ListSessionsByUserID(ctx context.Context, userID string) ([]authmodel.Session, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, token_hash, last_seen_at, expires_at, revoked_at, created_at, updated_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}
	defer rows.Close()

	result := make([]authmodel.Session, 0)
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历会话列表失败: %w", err)
	}
	return result, nil
}

func (s *SQLStore) RevokeSession(ctx context.Context, sessionID string, userID string, revokedAt time.Time, events []authmodel.AuthEvent) error {
	return ppdb.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		var returnedID string
		err := tx.QueryRowContext(ctx, `
			UPDATE user_sessions
			SET revoked_at = $3, updated_at = $3
			WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL
			RETURNING id
		`, sessionID, userID, revokedAt).Scan(&returnedID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return authmodel.ErrSessionNotFound
			}
			return fmt.Errorf("撤销会话失败: %w", err)
		}
		return insertAuthEventsTx(ctx, tx, events)
	})
}

func (s *SQLStore) FindSessionByTokenHash(ctx context.Context, tokenHash string) (authmodel.Session, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, last_seen_at, expires_at, revoked_at, created_at, updated_at
		FROM user_sessions
		WHERE token_hash = $1
	`, tokenHash)
	return scanSession(row)
}

func (s *SQLStore) ListPermissionsByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT permissions.code
		FROM user_roles
		JOIN role_permissions ON role_permissions.role_id = user_roles.role_id
		JOIN permissions ON permissions.id = role_permissions.permission_id
		WHERE user_roles.user_id = $1
		ORDER BY permissions.code ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("查询权限列表失败: %w", err)
	}
	defer rows.Close()
	return scanStringRows(rows)
}

func (s *SQLStore) ListRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT roles.code
		FROM user_roles
		JOIN roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = $1
		ORDER BY roles.code ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("查询角色列表失败: %w", err)
	}
	defer rows.Close()
	return scanStringRows(rows)
}

func consumePasswordResetToken(ctx context.Context, tx *sql.Tx, tokenHash string, usedAt time.Time) (authmodel.VerificationToken, error) {
	row := tx.QueryRowContext(ctx, `
		UPDATE verification_tokens
		SET used_at = $2, status = $3, updated_at = $2
		WHERE token_hash = $1
			AND purpose = $4
			AND used_at IS NULL
			AND expires_at >= $2
		RETURNING id, user_id, purpose, token_hash, destination, expires_at, used_at, status, created_at, updated_at
	`, tokenHash, usedAt, authmodel.VerificationTokenStatusConsumed, authmodel.VerificationPurposePasswordReset)
	token, err := scanVerificationToken(row)
	if err == nil {
		return token, nil
	}
	if !errors.Is(err, authmodel.ErrVerificationTokenNotFound) {
		return authmodel.VerificationToken{}, err
	}

	currentRow := tx.QueryRowContext(ctx, `
		SELECT id, user_id, purpose, token_hash, destination, expires_at, used_at, status, created_at, updated_at
		FROM verification_tokens
		WHERE token_hash = $1
		FOR UPDATE
	`, tokenHash)
	token, err = scanVerificationToken(currentRow)
	if err != nil {
		return authmodel.VerificationToken{}, err
	}
	switch {
	case token.Purpose != authmodel.VerificationPurposePasswordReset:
		return authmodel.VerificationToken{}, authmodel.ErrVerificationTokenPurposeMismatch
	case token.UsedAt != nil:
		return authmodel.VerificationToken{}, authmodel.ErrVerificationTokenConsumed
	case usedAt.After(token.ExpiresAt):
		return authmodel.VerificationToken{}, authmodel.ErrVerificationTokenExpired
	default:
		return authmodel.VerificationToken{}, authmodel.ErrVerificationTokenNotFound
	}
}

func insertAuthEventsTx(ctx context.Context, tx *sql.Tx, events []authmodel.AuthEvent) error {
	for _, event := range events {
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("序列化认证事件失败: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO auth_events(id, user_id, session_id, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		`, event.ID, nullString(event.UserID), nullString(event.SessionID), event.Type, string(payload), event.CreatedAt); err != nil {
			return fmt.Errorf("写入认证事件失败: %w", err)
		}
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (authmodel.User, error) {
	var user authmodel.User
	var archivedAt sql.NullTime
	if err := row.Scan(&user.ID, &user.Status, &archivedAt, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authmodel.User{}, authmodel.ErrUserNotFound
		}
		return authmodel.User{}, fmt.Errorf("扫描用户失败: %w", err)
	}
	if archivedAt.Valid {
		user.ArchivedAt = &archivedAt.Time
	}
	return user, nil
}

func scanIdentity(row rowScanner) (authmodel.Identity, error) {
	var identity authmodel.Identity
	var verifiedAt sql.NullTime
	if err := row.Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.Identifier,
		&identity.SecretHash,
		&identity.Status,
		&verifiedAt,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authmodel.Identity{}, authmodel.ErrIdentityNotFound
		}
		return authmodel.Identity{}, fmt.Errorf("扫描身份失败: %w", err)
	}
	if verifiedAt.Valid {
		identity.VerifiedAt = &verifiedAt.Time
	}
	return identity, nil
}

func scanSession(row rowScanner) (authmodel.Session, error) {
	var session authmodel.Session
	var revokedAt sql.NullTime
	if err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.LastSeenAt,
		&session.ExpiresAt,
		&revokedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authmodel.Session{}, authmodel.ErrSessionNotFound
		}
		return authmodel.Session{}, fmt.Errorf("扫描会话失败: %w", err)
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}
	return session, nil
}

func scanVerificationToken(row rowScanner) (authmodel.VerificationToken, error) {
	var token authmodel.VerificationToken
	var userID sql.NullString
	var usedAt sql.NullTime
	if err := row.Scan(
		&token.ID,
		&userID,
		&token.Purpose,
		&token.TokenHash,
		&token.Destination,
		&token.ExpiresAt,
		&usedAt,
		&token.Status,
		&token.CreatedAt,
		&token.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authmodel.VerificationToken{}, authmodel.ErrVerificationTokenNotFound
		}
		return authmodel.VerificationToken{}, fmt.Errorf("扫描验证码失败: %w", err)
	}
	if userID.Valid {
		token.UserID = userID.String
	}
	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}
	return token, nil
}

func scanStringRows(rows *sql.Rows) ([]string, error) {
	result := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("扫描字符串行失败: %w", err)
		}
		result = append(result, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历字符串行失败: %w", err)
	}
	return result, nil
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
