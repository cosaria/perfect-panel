package model

import (
	"context"
	"errors"
	"time"
)

const (
	IdentityProviderEmail = "email"

	UserStatusActive = "active"

	IdentityStatusActive = "active"

	VerificationPurposeEmailVerification = "email_verification"
	VerificationPurposePasswordReset     = "password_reset"

	VerificationTokenStatusPending  = "pending"
	VerificationTokenStatusConsumed = "consumed"

	PasswordResetStatusCompleted = "completed"

	AuthEventTypeSessionSignedIn        = "session_signed_in"
	AuthEventTypeSessionSignedOut       = "session_signed_out"
	AuthEventTypePasswordResetRequested = "password_reset_requested"
	AuthEventTypePasswordResetCompleted = "password_reset_completed"
	AuthEventTypeVerificationIssued     = "verification_issued"
)

var (
	ErrInvalidCredentials               = errors.New("认证凭据无效")
	ErrIdentityNotFound                 = errors.New("身份不存在")
	ErrUserNotFound                     = errors.New("用户不存在")
	ErrSessionNotFound                  = errors.New("会话不存在")
	ErrVerificationTokenNotFound        = errors.New("验证码不存在")
	ErrVerificationTokenConsumed        = errors.New("验证码已消费")
	ErrVerificationTokenExpired         = errors.New("验证码已过期")
	ErrVerificationTokenPurposeMismatch = errors.New("验证码用途不匹配")
	ErrUnauthorized                     = errors.New("未认证")
	ErrForbidden                        = errors.New("无权限")
)

type User struct {
	ID         string
	Status     string
	ArchivedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Identity struct {
	ID         string
	UserID     string
	Provider   string
	Identifier string
	SecretHash string
	Status     string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Session struct {
	ID         string
	UserID     string
	TokenHash  string
	LastSeenAt time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type VerificationToken struct {
	ID          string
	UserID      string
	Purpose     string
	TokenHash   string
	Destination string
	Status      string
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AuthEvent struct {
	ID        string
	UserID    string
	SessionID string
	Type      string
	Payload   map[string]any
	CreatedAt time.Time
}

type PasswordResetRecord struct {
	Token    VerificationToken
	Identity Identity
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResult struct {
	Session     Session
	AccessToken string
}

type IssueVerificationInput struct {
	Email   string `json:"email"`
	Purpose string
}

type IssuedVerificationToken struct {
	Token      VerificationToken
	PlainToken string
}

type RequestPasswordResetInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"password"`
}

type ResetPasswordResult struct {
	TokenID     string
	Status      string
	CompletedAt time.Time
}

type SignOutInput struct {
	UserID    string
	SessionID string
}

type Clock interface {
	Now() time.Time
}

type TokenGenerator interface {
	NewToken() (string, error)
}

type TokenHasher interface {
	Hash(token string) string
}

type PasswordHasher interface {
	Hash(secret string) (string, error)
	Compare(hash string, secret string) error
}

type SessionResolver interface {
	ResolveSession(ctx context.Context, bearerToken string) (Principal, error)
}

type Store interface {
	FindIdentityByIdentifier(ctx context.Context, provider string, identifier string) (Identity, error)
	FindUserByID(ctx context.Context, userID string) (User, error)
	CreateSession(ctx context.Context, session Session, events []AuthEvent) error
	SaveVerificationToken(ctx context.Context, token VerificationToken, events []AuthEvent) error
	ConsumePasswordReset(ctx context.Context, tokenHash string, passwordHash string, usedAt time.Time, events []AuthEvent) (PasswordResetRecord, error)
	ListSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
	RevokeSession(ctx context.Context, sessionID string, userID string, revokedAt time.Time, events []AuthEvent) error
	FindSessionByTokenHash(ctx context.Context, tokenHash string) (Session, error)
	ListPermissionsByUserID(ctx context.Context, userID string) ([]string, error)
	ListRolesByUserID(ctx context.Context, userID string) ([]string, error)
}

type RequiredSetting struct {
	Scope string
	Key   string
	Value any
}

func RequiredSettings() []RequiredSetting {
	return []RequiredSetting{
		{Scope: "site", Key: "site_name", Value: "Perfect Panel"},
		{Scope: "site", Key: "app_name", Value: "Perfect Panel"},
		{Scope: "auth", Key: "session_ttl_seconds", Value: 86400},
		{Scope: "auth", Key: "verification_token_ttl_seconds", Value: 1800},
		{Scope: "auth", Key: "password_reset_token_ttl_seconds", Value: 1800},
		{Scope: "auth", Key: "email_identity_enabled", Value: true},
	}
}

type Principal struct {
	UserID      string   `json:"userId"`
	SessionID   string   `json:"sessionId"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type principalContextKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok
}
