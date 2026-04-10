package auth

import (
	"context"
	"database/sql"
	"net/http"

	authapi "github.com/perfect-panel/server-v2/internal/domains/auth/api"
	authmodel "github.com/perfect-panel/server-v2/internal/domains/auth/model"
	authstore "github.com/perfect-panel/server-v2/internal/domains/auth/store"
	authusecase "github.com/perfect-panel/server-v2/internal/domains/auth/usecase"
)

const (
	IdentityProviderEmail = authmodel.IdentityProviderEmail

	UserStatusActive = authmodel.UserStatusActive

	IdentityStatusActive = authmodel.IdentityStatusActive

	VerificationPurposeEmailVerification = authmodel.VerificationPurposeEmailVerification
	VerificationPurposePasswordReset     = authmodel.VerificationPurposePasswordReset

	VerificationTokenStatusPending  = authmodel.VerificationTokenStatusPending
	VerificationTokenStatusConsumed = authmodel.VerificationTokenStatusConsumed

	PasswordResetStatusCompleted = authmodel.PasswordResetStatusCompleted

	AuthEventTypeSessionSignedIn        = authmodel.AuthEventTypeSessionSignedIn
	AuthEventTypeSessionSignedOut       = authmodel.AuthEventTypeSessionSignedOut
	AuthEventTypePasswordResetRequested = authmodel.AuthEventTypePasswordResetRequested
	AuthEventTypePasswordResetCompleted = authmodel.AuthEventTypePasswordResetCompleted
	AuthEventTypeVerificationIssued     = authmodel.AuthEventTypeVerificationIssued
)

var (
	ErrInvalidCredentials               = authmodel.ErrInvalidCredentials
	ErrIdentityNotFound                 = authmodel.ErrIdentityNotFound
	ErrUserNotFound                     = authmodel.ErrUserNotFound
	ErrSessionNotFound                  = authmodel.ErrSessionNotFound
	ErrVerificationTokenNotFound        = authmodel.ErrVerificationTokenNotFound
	ErrVerificationTokenConsumed        = authmodel.ErrVerificationTokenConsumed
	ErrVerificationTokenExpired         = authmodel.ErrVerificationTokenExpired
	ErrVerificationTokenPurposeMismatch = authmodel.ErrVerificationTokenPurposeMismatch
	ErrUnauthorized                     = authmodel.ErrUnauthorized
	ErrForbidden                        = authmodel.ErrForbidden
)

type User = authmodel.User
type Identity = authmodel.Identity
type Session = authmodel.Session
type VerificationToken = authmodel.VerificationToken
type AuthEvent = authmodel.AuthEvent
type PasswordResetRecord = authmodel.PasswordResetRecord
type SignInInput = authmodel.SignInInput
type SignInResult = authmodel.SignInResult
type IssueVerificationInput = authmodel.IssueVerificationInput
type IssuedVerificationToken = authmodel.IssuedVerificationToken
type RequestPasswordResetInput = authmodel.RequestPasswordResetInput
type ResetPasswordInput = authmodel.ResetPasswordInput
type ResetPasswordResult = authmodel.ResetPasswordResult
type SignOutInput = authmodel.SignOutInput
type Clock = authmodel.Clock
type TokenGenerator = authmodel.TokenGenerator
type TokenHasher = authmodel.TokenHasher
type PasswordHasher = authmodel.PasswordHasher
type SessionResolver = authmodel.SessionResolver
type Store = authmodel.Store
type RequiredSetting = authmodel.RequiredSetting
type Principal = authmodel.Principal

type ServiceOptions = authusecase.ServiceOptions
type Service = authusecase.Service
type DefaultSessionResolver = authusecase.DefaultSessionResolver
type HTTPHandler = authapi.HTTPHandler
type SQLStore = authstore.SQLStore

func RequiredSettings() []RequiredSetting {
	return authmodel.RequiredSettings()
}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return authmodel.WithPrincipal(ctx, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	return authmodel.PrincipalFromContext(ctx)
}

func NewService(opts ServiceOptions) *Service {
	return authusecase.NewService(opts)
}

func NewHTTPHandler(service *Service) *HTTPHandler {
	return authapi.NewHTTPHandler(service)
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return authstore.NewSQLStore(db)
}

func WriteEnvelope(w http.ResponseWriter, status int, data any) {
	authapi.WriteEnvelope(w, status, data)
}

func WriteEnvelopeWithMeta(w http.ResponseWriter, status int, data any, meta map[string]any) {
	authapi.WriteEnvelopeWithMeta(w, status, data, meta)
}

func WriteProblem(w http.ResponseWriter, status int, title string, detail string, code string) {
	authapi.WriteProblem(w, status, title, detail, code)
}
