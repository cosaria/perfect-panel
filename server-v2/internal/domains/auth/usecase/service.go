package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	authmodel "github.com/perfect-panel/server-v2/internal/domains/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type ServiceOptions struct {
	Store           authmodel.Store
	PasswordHasher  authmodel.PasswordHasher
	TokenHasher     authmodel.TokenHasher
	Clock           authmodel.Clock
	SessionTTL      time.Duration
	TokenTTL        time.Duration
	TokenGenerator  authmodel.TokenGenerator
	SessionResolver authmodel.SessionResolver
}

type Service struct {
	store          authmodel.Store
	passwordHasher authmodel.PasswordHasher
	tokenHasher    authmodel.TokenHasher
	clock          authmodel.Clock
	sessionTTL     time.Duration
	tokenTTL       time.Duration
	tokenGenerator authmodel.TokenGenerator
}

func NewService(opts ServiceOptions) *Service {
	return &Service{
		store:          opts.Store,
		passwordHasher: coalescePasswordHasher(opts.PasswordHasher),
		tokenHasher:    coalesceTokenHasher(opts.TokenHasher),
		clock:          coalesceClock(opts.Clock),
		sessionTTL:     coalesceDuration(opts.SessionTTL, 24*time.Hour),
		tokenTTL:       coalesceDuration(opts.TokenTTL, 30*time.Minute),
		tokenGenerator: coalesceTokenGenerator(opts.TokenGenerator),
	}
}

func (s *Service) Clock() authmodel.Clock {
	return s.clock
}

func (s *Service) SignIn(ctx context.Context, input authmodel.SignInInput) (authmodel.SignInResult, error) {
	identifier := normalizeEmail(input.Email)
	if identifier == "" || strings.TrimSpace(input.Password) == "" {
		return authmodel.SignInResult{}, authmodel.ErrInvalidCredentials
	}

	identity, err := s.store.FindIdentityByIdentifier(ctx, authmodel.IdentityProviderEmail, identifier)
	if err != nil {
		return authmodel.SignInResult{}, err
	}
	user, err := s.store.FindUserByID(ctx, identity.UserID)
	if err != nil {
		return authmodel.SignInResult{}, err
	}
	if user.Status != "" && user.Status != authmodel.UserStatusActive {
		return authmodel.SignInResult{}, authmodel.ErrUnauthorized
	}
	if err := s.passwordHasher.Compare(identity.SecretHash, input.Password); err != nil {
		return authmodel.SignInResult{}, err
	}

	now := s.clock.Now().UTC()
	accessToken, err := s.tokenGenerator.NewToken()
	if err != nil {
		return authmodel.SignInResult{}, fmt.Errorf("生成会话 token 失败: %w", err)
	}
	sessionID, err := newUUID()
	if err != nil {
		return authmodel.SignInResult{}, err
	}

	session := authmodel.Session{
		ID:         sessionID,
		UserID:     user.ID,
		TokenHash:  s.tokenHasher.Hash(accessToken),
		LastSeenAt: now,
		ExpiresAt:  now.Add(s.sessionTTL),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	events, err := s.prepareEvents(authmodel.AuthEvent{
		UserID:    user.ID,
		SessionID: session.ID,
		Type:      authmodel.AuthEventTypeSessionSignedIn,
		Payload: map[string]any{
			"provider": authmodel.IdentityProviderEmail,
		},
		CreatedAt: now,
	})
	if err != nil {
		return authmodel.SignInResult{}, err
	}
	if err := s.store.CreateSession(ctx, session, events); err != nil {
		return authmodel.SignInResult{}, err
	}

	return authmodel.SignInResult{
		Session:     session,
		AccessToken: accessToken,
	}, nil
}

func (s *Service) IssueVerificationToken(ctx context.Context, input authmodel.IssueVerificationInput) (authmodel.IssuedVerificationToken, error) {
	identifier := normalizeEmail(input.Email)
	if identifier == "" {
		return authmodel.IssuedVerificationToken{}, authmodel.ErrIdentityNotFound
	}
	if strings.TrimSpace(input.Purpose) == "" {
		input.Purpose = authmodel.VerificationPurposeEmailVerification
	}

	identity, err := s.store.FindIdentityByIdentifier(ctx, authmodel.IdentityProviderEmail, identifier)
	if err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}

	return s.issueVerificationToken(ctx, identity.UserID, identifier, input.Purpose)
}

func (s *Service) RequestPasswordReset(ctx context.Context, input authmodel.RequestPasswordResetInput) (authmodel.IssuedVerificationToken, error) {
	identifier := normalizeEmail(input.Email)
	if identifier == "" {
		return authmodel.IssuedVerificationToken{}, authmodel.ErrIdentityNotFound
	}

	identity, err := s.store.FindIdentityByIdentifier(ctx, authmodel.IdentityProviderEmail, identifier)
	if err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}

	issued, err := s.issueVerificationToken(ctx, identity.UserID, identifier, authmodel.VerificationPurposePasswordReset)
	if err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}

	return issued, nil
}

func (s *Service) ResetPassword(ctx context.Context, input authmodel.ResetPasswordInput) (authmodel.ResetPasswordResult, error) {
	if strings.TrimSpace(input.Token) == "" || strings.TrimSpace(input.NewPassword) == "" {
		return authmodel.ResetPasswordResult{}, authmodel.ErrVerificationTokenNotFound
	}

	now := s.clock.Now().UTC()
	passwordHash, err := s.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return authmodel.ResetPasswordResult{}, fmt.Errorf("计算密码哈希失败: %w", err)
	}
	events, err := s.prepareEvents(authmodel.AuthEvent{
		Type: authmodel.AuthEventTypePasswordResetCompleted,
		Payload: map[string]any{
			"tokenHash": s.tokenHasher.Hash(input.Token),
		},
		CreatedAt: now,
	})
	if err != nil {
		return authmodel.ResetPasswordResult{}, err
	}
	record, err := s.store.ConsumePasswordReset(ctx, s.tokenHasher.Hash(input.Token), passwordHash, now, events)
	if err != nil {
		return authmodel.ResetPasswordResult{}, err
	}

	return authmodel.ResetPasswordResult{
		TokenID:     record.Token.ID,
		Status:      authmodel.PasswordResetStatusCompleted,
		CompletedAt: now,
	}, nil
}

func (s *Service) SignOut(ctx context.Context, input authmodel.SignOutInput) error {
	now := s.clock.Now().UTC()
	events, err := s.prepareEvents(authmodel.AuthEvent{
		UserID:    input.UserID,
		SessionID: input.SessionID,
		Type:      authmodel.AuthEventTypeSessionSignedOut,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}
	return s.store.RevokeSession(ctx, input.SessionID, input.UserID, now, events)
}

func (s *Service) ListUserSessions(ctx context.Context, userID string) ([]authmodel.Session, error) {
	return s.store.ListSessionsByUserID(ctx, userID)
}

func (s *Service) issueVerificationToken(ctx context.Context, userID string, destination string, purpose string) (authmodel.IssuedVerificationToken, error) {
	now := s.clock.Now().UTC()
	rawToken, err := s.tokenGenerator.NewToken()
	if err != nil {
		return authmodel.IssuedVerificationToken{}, fmt.Errorf("生成验证码失败: %w", err)
	}
	tokenID, err := newUUID()
	if err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}

	token := authmodel.VerificationToken{
		ID:          tokenID,
		UserID:      userID,
		Purpose:     purpose,
		TokenHash:   s.tokenHasher.Hash(rawToken),
		Destination: destination,
		Status:      authmodel.VerificationTokenStatusPending,
		ExpiresAt:   now.Add(s.tokenTTL),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	events := []authmodel.AuthEvent{
		{
			UserID: userID,
			Type:   authmodel.AuthEventTypeVerificationIssued,
			Payload: map[string]any{
				"tokenId":     token.ID,
				"destination": destination,
				"purpose":     purpose,
			},
			CreatedAt: now,
		},
	}
	if purpose == authmodel.VerificationPurposePasswordReset {
		events = append(events, authmodel.AuthEvent{
			UserID: userID,
			Type:   authmodel.AuthEventTypePasswordResetRequested,
			Payload: map[string]any{
				"destination": destination,
				"purpose":     purpose,
			},
			CreatedAt: now,
		})
	}
	preparedEvents, err := s.prepareEvents(events...)
	if err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}
	if err := s.store.SaveVerificationToken(ctx, token, preparedEvents); err != nil {
		return authmodel.IssuedVerificationToken{}, err
	}
	return authmodel.IssuedVerificationToken{
		Token:      token,
		PlainToken: rawToken,
	}, nil
}

func (s *Service) prepareEvents(events ...authmodel.AuthEvent) ([]authmodel.AuthEvent, error) {
	prepared := make([]authmodel.AuthEvent, 0, len(events))
	for _, event := range events {
		if event.ID == "" {
			eventID, err := newUUID()
			if err != nil {
				return nil, err
			}
			event.ID = eventID
		}
		if event.Payload == nil {
			event.Payload = map[string]any{}
		}
		prepared = append(prepared, event)
	}
	return prepared, nil
}

type DefaultSessionResolver struct {
	Store       authmodel.Store
	TokenHasher authmodel.TokenHasher
	Clock       authmodel.Clock
}

func (r DefaultSessionResolver) ResolveSession(ctx context.Context, bearerToken string) (authmodel.Principal, error) {
	if r.Store == nil {
		return authmodel.Principal{}, authmodel.ErrUnauthorized
	}
	tokenHash := coalesceTokenHasher(r.TokenHasher).Hash(strings.TrimSpace(bearerToken))
	session, err := r.Store.FindSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return authmodel.Principal{}, err
	}
	now := coalesceClock(r.Clock).Now().UTC()
	if session.RevokedAt != nil || now.After(session.ExpiresAt) {
		return authmodel.Principal{}, authmodel.ErrUnauthorized
	}
	user, err := r.Store.FindUserByID(ctx, session.UserID)
	if err != nil {
		return authmodel.Principal{}, err
	}
	if user.Status != "" && user.Status != authmodel.UserStatusActive {
		return authmodel.Principal{}, authmodel.ErrUnauthorized
	}
	permissions, err := r.Store.ListPermissionsByUserID(ctx, session.UserID)
	if err != nil {
		return authmodel.Principal{}, err
	}
	roles, err := r.Store.ListRolesByUserID(ctx, session.UserID)
	if err != nil {
		return authmodel.Principal{}, err
	}
	return authmodel.Principal{
		UserID:      session.UserID,
		SessionID:   session.ID,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

type bcryptHasher struct{}

func (bcryptHasher) Hash(secret string) (string, error) {
	raw, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (bcryptHasher) Compare(hash string, secret string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret)); err != nil {
		return authmodel.ErrInvalidCredentials
	}
	return nil
}

type sha256Hasher struct{}

func (sha256Hasher) Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type secureTokenGenerator struct{}

func (secureTokenGenerator) NewToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("生成随机 token 失败: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func coalesceClock(clock authmodel.Clock) authmodel.Clock {
	if clock != nil {
		return clock
	}
	return systemClock{}
}

func coalesceTokenHasher(hasher authmodel.TokenHasher) authmodel.TokenHasher {
	if hasher != nil {
		return hasher
	}
	return sha256Hasher{}
}

func coalescePasswordHasher(hasher authmodel.PasswordHasher) authmodel.PasswordHasher {
	if hasher != nil {
		return hasher
	}
	return bcryptHasher{}
}

func coalesceTokenGenerator(generator authmodel.TokenGenerator) authmodel.TokenGenerator {
	if generator != nil {
		return generator
	}
	return secureTokenGenerator{}
}

func coalesceDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return fallback
}

func newUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("生成 UUID 失败: %w", err)
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16],
	), nil
}
