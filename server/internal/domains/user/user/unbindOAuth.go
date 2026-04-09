package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type UnbindOAuthInput struct {
	Body types.UnbindOAuthRequest
}

func UnbindOAuthHandler(deps Deps) func(context.Context, *UnbindOAuthInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnbindOAuthInput) (*struct{}, error) {
		l := NewUnbindOAuthLogic(ctx, deps)
		if err := l.UnbindOAuth(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UnbindOAuthLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Unbind OAuth
func NewUnbindOAuthLogic(ctx context.Context, deps Deps) *UnbindOAuthLogic {
	return &UnbindOAuthLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UnbindOAuthLogic) UnbindOAuth(req *types.UnbindOAuthRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	if !l.validator(req) {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "invalid parameter")
	}
	err := l.deps.UserModel.DeleteUserAuthMethods(l.ctx, u.Id, req.Method)
	if err != nil {
		l.Errorw("delete user auth methods failed:", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete user auth methods failed: %v", err.Error())
	}

	return nil
}
func (l *UnbindOAuthLogic) validator(req *types.UnbindOAuthRequest) bool {
	return req.Method != "" && req.Method != "email" && req.Method != "mobile"
}
