package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CurrentUserOutput struct {
	Body *types.User
}

func CurrentUserHandler(deps Deps) func(context.Context, *struct{}) (*CurrentUserOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*CurrentUserOutput, error) {
		l := NewCurrentUserLogic(ctx, deps)
		resp, err := l.CurrentUser()
		if err != nil {
			return nil, err
		}
		return &CurrentUserOutput{Body: resp}, nil
	}
}

type CurrentUserLogic struct {
	ctx  context.Context
	deps Deps
	logger.Logger
}

func NewCurrentUserLogic(ctx context.Context, deps Deps) *CurrentUserLogic {
	return &CurrentUserLogic{
		ctx:    ctx,
		deps:   deps,
		Logger: logger.WithContext(ctx),
	}
}

func (l *CurrentUserLogic) CurrentUser() (*types.User, error) {
	resp := &types.User{}
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	l.Info("current user", zap.Field{Key: "userId", Type: zapcore.Int64Type, Integer: u.Id})
	tool.DeepCopy(resp, u)
	return resp, nil
}
