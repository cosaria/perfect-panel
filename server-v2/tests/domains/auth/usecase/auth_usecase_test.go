package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
)

func TestSignInCreatesSession(t *testing.T) {
	t.Parallel()

	clock := time.Date(2026, 4, 10, 9, 0, 0, 0, time.UTC)
	store := newAuthStoreStub()
	store.users["user-1"] = authdomain.User{ID: "user-1", Status: authdomain.UserStatusActive}
	store.identities[identityKey(authdomain.IdentityProviderEmail, "admin@example.com")] = authdomain.Identity{
		ID:         "identity-1",
		UserID:     "user-1",
		Provider:   authdomain.IdentityProviderEmail,
		Identifier: "admin@example.com",
		SecretHash: "hashed:password-123",
		Status:     authdomain.IdentityStatusActive,
	}

	service := authdomain.NewService(authdomain.ServiceOptions{
		Store:           store,
		PasswordHasher:  passwordHasherStub{},
		TokenHasher:     tokenHasherStub{},
		Clock:           fixedClock(clock),
		SessionTTL:      2 * time.Hour,
		TokenGenerator:  fixedTokenGenerator("session-plain"),
		SessionResolver: nil,
	})

	result, err := service.SignIn(context.Background(), authdomain.SignInInput{
		Email:    "admin@example.com",
		Password: "password-123",
	})
	if err != nil {
		t.Fatalf("SignIn 返回错误: %v", err)
	}

	if result.Session.UserID != "user-1" {
		t.Fatalf("会话 user_id 不匹配: got=%s", result.Session.UserID)
	}
	if result.AccessToken != "session-plain" {
		t.Fatalf("返回的 access token 不匹配: got=%s", result.AccessToken)
	}
	if len(store.createdSessions) != 1 {
		t.Fatalf("应创建 1 条会话，实际 %d", len(store.createdSessions))
	}
	if got := store.createdSessions[0].TokenHash; got != "token:session-plain" {
		t.Fatalf("会话应保存 token 哈希，got=%s", got)
	}
	if store.createdSessions[0].TokenHash == result.AccessToken {
		t.Fatal("会话不应保存明文 token")
	}
	if want := clock.Add(2 * time.Hour); !result.Session.ExpiresAt.Equal(want) {
		t.Fatalf("会话过期时间不匹配: want=%s got=%s", want, result.Session.ExpiresAt)
	}
	if len(store.authEvents) != 1 || store.authEvents[0].Type != authdomain.AuthEventTypeSessionSignedIn {
		t.Fatalf("应写入 1 条登录事件，实际 %+v", store.authEvents)
	}
}

func TestRequestPasswordResetStoresHashedTokenOnly(t *testing.T) {
	t.Parallel()

	clock := time.Date(2026, 4, 10, 9, 30, 0, 0, time.UTC)
	store := newAuthStoreStub()
	store.users["user-1"] = authdomain.User{ID: "user-1", Status: authdomain.UserStatusActive}
	store.identities[identityKey(authdomain.IdentityProviderEmail, "user@example.com")] = authdomain.Identity{
		ID:         "identity-1",
		UserID:     "user-1",
		Provider:   authdomain.IdentityProviderEmail,
		Identifier: "user@example.com",
		SecretHash: "hashed:old-password",
		Status:     authdomain.IdentityStatusActive,
	}

	service := authdomain.NewService(authdomain.ServiceOptions{
		Store:          store,
		PasswordHasher: passwordHasherStub{},
		TokenHasher:    tokenHasherStub{},
		Clock:          fixedClock(clock),
		TokenTTL:       30 * time.Minute,
		TokenGenerator: fixedTokenGenerator("reset-plain"),
	})

	result, err := service.RequestPasswordReset(context.Background(), authdomain.RequestPasswordResetInput{
		Email: "user@example.com",
	})
	if err != nil {
		t.Fatalf("RequestPasswordReset 返回错误: %v", err)
	}

	if result.PlainToken != "reset-plain" {
		t.Fatalf("返回的明文 token 不匹配: got=%s", result.PlainToken)
	}
	if len(store.savedVerificationTokens) != 1 {
		t.Fatalf("应保存 1 条验证码记录，实际 %d", len(store.savedVerificationTokens))
	}
	token := store.savedVerificationTokens[0]
	if token.TokenHash != "token:reset-plain" {
		t.Fatalf("验证码应保存哈希，got=%s", token.TokenHash)
	}
	if token.TokenHash == result.PlainToken {
		t.Fatal("验证码不应保存明文 token")
	}
	if token.Purpose != authdomain.VerificationPurposePasswordReset {
		t.Fatalf("验证码用途不匹配: got=%s", token.Purpose)
	}
	if want := clock.Add(30 * time.Minute); !token.ExpiresAt.Equal(want) {
		t.Fatalf("验证码过期时间不匹配: want=%s got=%s", want, token.ExpiresAt)
	}
}

func TestResetPasswordConsumesVerificationTokenOnce(t *testing.T) {
	t.Parallel()

	clock := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	store := newAuthStoreStub()
	store.users["user-1"] = authdomain.User{ID: "user-1", Status: authdomain.UserStatusActive}
	store.identities[identityKey(authdomain.IdentityProviderEmail, "user@example.com")] = authdomain.Identity{
		ID:         "identity-1",
		UserID:     "user-1",
		Provider:   authdomain.IdentityProviderEmail,
		Identifier: "user@example.com",
		SecretHash: "hashed:old-password",
		Status:     authdomain.IdentityStatusActive,
	}
	store.verificationByHash["token:reset-plain"] = authdomain.VerificationToken{
		ID:          "verification-1",
		UserID:      "user-1",
		Purpose:     authdomain.VerificationPurposePasswordReset,
		TokenHash:   "token:reset-plain",
		Destination: "user@example.com",
		Status:      authdomain.VerificationTokenStatusPending,
		ExpiresAt:   clock.Add(15 * time.Minute),
	}

	service := authdomain.NewService(authdomain.ServiceOptions{
		Store:          store,
		PasswordHasher: passwordHasherStub{},
		TokenHasher:    tokenHasherStub{},
		Clock:          fixedClock(clock),
	})

	result, err := service.ResetPassword(context.Background(), authdomain.ResetPasswordInput{
		Token:       "reset-plain",
		NewPassword: "new-password-123",
	})
	if err != nil {
		t.Fatalf("ResetPassword 第一次执行返回错误: %v", err)
	}
	if result.Status != authdomain.PasswordResetStatusCompleted {
		t.Fatalf("密码重置状态不匹配: got=%s", result.Status)
	}
	updatedIdentity := store.identities[identityKey(authdomain.IdentityProviderEmail, "user@example.com")]
	if updatedIdentity.SecretHash != "hashed:new-password-123" {
		t.Fatalf("密码哈希未更新: got=%s", updatedIdentity.SecretHash)
	}
	usedToken := store.verificationByHash["token:reset-plain"]
	if usedToken.UsedAt == nil || !usedToken.UsedAt.Equal(clock) {
		t.Fatalf("验证码未被标记为消费: %+v", usedToken.UsedAt)
	}

	_, err = service.ResetPassword(context.Background(), authdomain.ResetPasswordInput{
		Token:       "reset-plain",
		NewPassword: "another-password-123",
	})
	if !errors.Is(err, authdomain.ErrVerificationTokenConsumed) {
		t.Fatalf("第二次消费同一 token 应返回 ErrVerificationTokenConsumed，实际 %v", err)
	}
}

type authStoreStub struct {
	users                   map[string]authdomain.User
	identities              map[string]authdomain.Identity
	verificationByHash      map[string]authdomain.VerificationToken
	createdSessions         []authdomain.Session
	savedVerificationTokens []authdomain.VerificationToken
	authEvents              []authdomain.AuthEvent
	roleCodes               map[string][]string
	permissionCodes         map[string][]string
}

func newAuthStoreStub() *authStoreStub {
	return &authStoreStub{
		users:              make(map[string]authdomain.User),
		identities:         make(map[string]authdomain.Identity),
		verificationByHash: make(map[string]authdomain.VerificationToken),
		roleCodes:          make(map[string][]string),
		permissionCodes:    make(map[string][]string),
	}
}

func (s *authStoreStub) FindIdentityByIdentifier(_ context.Context, provider string, identifier string) (authdomain.Identity, error) {
	identity, ok := s.identities[identityKey(provider, identifier)]
	if !ok {
		return authdomain.Identity{}, authdomain.ErrIdentityNotFound
	}
	return identity, nil
}

func (s *authStoreStub) FindUserByID(_ context.Context, userID string) (authdomain.User, error) {
	user, ok := s.users[userID]
	if !ok {
		return authdomain.User{}, authdomain.ErrUserNotFound
	}
	return user, nil
}

func (s *authStoreStub) CreateSession(_ context.Context, session authdomain.Session, events []authdomain.AuthEvent) error {
	s.createdSessions = append(s.createdSessions, session)
	s.authEvents = append(s.authEvents, events...)
	return nil
}

func (s *authStoreStub) SaveVerificationToken(_ context.Context, token authdomain.VerificationToken, events []authdomain.AuthEvent) error {
	s.savedVerificationTokens = append(s.savedVerificationTokens, token)
	s.verificationByHash[token.TokenHash] = token
	s.authEvents = append(s.authEvents, events...)
	return nil
}

func (s *authStoreStub) ConsumePasswordReset(_ context.Context, tokenHash string, passwordHash string, usedAt time.Time, events []authdomain.AuthEvent) (authdomain.PasswordResetRecord, error) {
	token, ok := s.verificationByHash[tokenHash]
	if !ok {
		return authdomain.PasswordResetRecord{}, authdomain.ErrVerificationTokenNotFound
	}
	if token.Purpose != authdomain.VerificationPurposePasswordReset {
		return authdomain.PasswordResetRecord{}, authdomain.ErrVerificationTokenPurposeMismatch
	}
	if token.UsedAt != nil {
		return authdomain.PasswordResetRecord{}, authdomain.ErrVerificationTokenConsumed
	}
	if usedAt.After(token.ExpiresAt) {
		return authdomain.PasswordResetRecord{}, authdomain.ErrVerificationTokenExpired
	}

	identity := s.identities[identityKey(authdomain.IdentityProviderEmail, token.Destination)]
	identity.SecretHash = passwordHash
	s.identities[identityKey(identity.Provider, identity.Identifier)] = identity

	token.Status = authdomain.VerificationTokenStatusConsumed
	token.UsedAt = &usedAt
	s.verificationByHash[tokenHash] = token
	enrichedEvents := append([]authdomain.AuthEvent(nil), events...)
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
	s.authEvents = append(s.authEvents, enrichedEvents...)

	return authdomain.PasswordResetRecord{
		Token:    token,
		Identity: identity,
	}, nil
}

func (s *authStoreStub) ListSessionsByUserID(_ context.Context, userID string) ([]authdomain.Session, error) {
	result := make([]authdomain.Session, 0, len(s.createdSessions))
	for _, session := range s.createdSessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}
	return result, nil
}

func (s *authStoreStub) RevokeSession(_ context.Context, sessionID string, userID string, revokedAt time.Time, events []authdomain.AuthEvent) error {
	for idx, session := range s.createdSessions {
		if session.ID == sessionID && session.UserID == userID {
			session.RevokedAt = &revokedAt
			s.createdSessions[idx] = session
			s.authEvents = append(s.authEvents, events...)
			return nil
		}
	}
	return authdomain.ErrSessionNotFound
}

func (s *authStoreStub) FindSessionByTokenHash(_ context.Context, tokenHash string) (authdomain.Session, error) {
	for _, session := range s.createdSessions {
		if session.TokenHash == tokenHash {
			return session, nil
		}
	}
	return authdomain.Session{}, authdomain.ErrSessionNotFound
}

func (s *authStoreStub) ListPermissionsByUserID(_ context.Context, userID string) ([]string, error) {
	return s.permissionCodes[userID], nil
}

func (s *authStoreStub) ListRolesByUserID(_ context.Context, userID string) ([]string, error) {
	return s.roleCodes[userID], nil
}

type passwordHasherStub struct{}

func (passwordHasherStub) Hash(secret string) (string, error) {
	return "hashed:" + secret, nil
}

func (passwordHasherStub) Compare(hash string, secret string) error {
	if hash != "hashed:"+secret {
		return authdomain.ErrInvalidCredentials
	}
	return nil
}

type tokenHasherStub struct{}

func (tokenHasherStub) Hash(token string) string {
	return "token:" + token
}

type fixedClock time.Time

func (c fixedClock) Now() time.Time {
	return time.Time(c)
}

type fixedTokenGenerator string

func (g fixedTokenGenerator) NewToken() (string, error) {
	return string(g), nil
}

func identityKey(provider string, identifier string) string {
	return provider + "::" + identifier
}

var _ authdomain.Store = (*authStoreStub)(nil)

func TestPasswordResetRecordShapeStaysStable(t *testing.T) {
	t.Parallel()

	recordType := reflect.TypeOf(authdomain.PasswordResetRecord{})
	if _, ok := recordType.FieldByName("Token"); !ok {
		t.Fatal("PasswordResetRecord 必须暴露 Token 字段")
	}
	if _, ok := recordType.FieldByName("Identity"); !ok {
		t.Fatal("PasswordResetRecord 必须暴露 Identity 字段")
	}
}
