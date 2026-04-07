package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"

	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
