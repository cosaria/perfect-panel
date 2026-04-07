package user

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CurrentUserOutput struct {
	Body *types.User
}

func CurrentUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*CurrentUserOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*CurrentUserOutput, error) {
		l := NewCurrentUserLogic(ctx, svcCtx)
		resp, err := l.CurrentUser()
		if err != nil {
			return nil, err
		}
		return &CurrentUserOutput{Body: resp}, nil
	}
}

type CurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger.Logger
}

func NewCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserLogic {
	return &CurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
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

	l.Logger.Info("current user", zap.Field{Key: "userId", Type: zapcore.Int64Type, Integer: u.Id})
	tool.DeepCopy(resp, u)
	return resp, nil
}
