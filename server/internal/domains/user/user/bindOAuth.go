package user

import (
	"context"
	"fmt"

	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/support/auth/oauth/google"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/random"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type BindOAuthInput struct {
	Body types.BindOAuthRequest
}

type BindOAuthOutput struct {
	Body *types.BindOAuthResponse
}

func BindOAuthHandler(deps Deps) func(context.Context, *BindOAuthInput) (*BindOAuthOutput, error) {
	return func(ctx context.Context, input *BindOAuthInput) (*BindOAuthOutput, error) {
		l := NewBindOAuthLogic(ctx, deps)
		resp, err := l.BindOAuth(&input.Body)
		if err != nil {
			return nil, err
		}
		return &BindOAuthOutput{Body: resp}, nil
	}
}

type BindOAuthLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Bind OAuth
func NewBindOAuthLogic(ctx context.Context, deps Deps) *BindOAuthLogic {
	return &BindOAuthLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *BindOAuthLogic) BindOAuth(req *types.BindOAuthRequest) (resp *types.BindOAuthResponse, err error) {
	var uri string
	switch req.Method {
	case "google":
		uri, err = l.google(req)
	case "apple":
		uri, err = l.apple(req)
	case "github":
		uri, err = l.github()
	case "facebook":
		uri, err = l.facebook()
	default:
		l.Errorw("oauth login method not support: %v", logger.Field("method", req.Method))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "oauth login method not support: %v", req.Method)
	}
	if err != nil {
		l.Errorw("error bind oauth", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "error bind oauth: %v", err.Error())
	}
	return &types.BindOAuthResponse{
		Redirect: uri,
	}, nil
}

func (l *BindOAuthLogic) google(req *types.BindOAuthRequest) (string, error) {
	authMethod, err := l.deps.AuthModel.FindOneByMethod(l.ctx, "google")
	if err != nil {
		return "", err
	}
	cfg := new(auth.GoogleAuthConfig)
	err = cfg.Unmarshal(authMethod.Config)
	if err != nil {
		l.Errorw("error unmarshal google config: %v", logger.Field("config", authMethod.Config), logger.Field("error", err.Error()))
		return "", err
	}
	client := google.New(&google.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  req.Redirect,
	})
	// generate the state code
	code := random.KeyNew(8, 1)
	// save the state code
	err = l.deps.Redis.Set(l.ctx, fmt.Sprintf("google:%s", code), req.Redirect, 5*60*time.Second).Err()
	if err != nil {
		return "", err
	}
	uri := client.AuthCodeURL(code, oauth2.AccessTypeOffline)
	return uri, nil
}

func (l *BindOAuthLogic) facebook() (string, error) {
	return "", nil
}
func (l *BindOAuthLogic) apple(req *types.BindOAuthRequest) (string, error) {
	authMethod, err := l.deps.AuthModel.FindOneByMethod(l.ctx, "apple")
	if err != nil {
		return "", err
	}
	var cfg auth.AppleAuthConfig
	err = cfg.Unmarshal(authMethod.Config)
	if err != nil {
		l.Errorw("error unmarshal apple config: %v", logger.Field("config", authMethod.Config), logger.Field("error", err.Error()))
		return "", err
	}
	uri := "https://appleid.apple.com/auth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=name email&response_mode=form_post"
	// generate the state code
	code := random.KeyNew(8, 1)
	// save the state code
	err = l.deps.Redis.Set(l.ctx, fmt.Sprintf("apple:%s", code), req.Redirect, 5*60*time.Second).Err()
	if err != nil {
		l.Errorw("error save state code to redis: %v", logger.Field("code", code), logger.Field("error", err.Error()))
	}
	return fmt.Sprintf(uri, cfg.ClientId, fmt.Sprintf("%s/api/v1/auth/oauth/callback/apple", cfg.RedirectURL), code), nil
}
func (l *BindOAuthLogic) github() (string, error) {
	return "", nil
}
