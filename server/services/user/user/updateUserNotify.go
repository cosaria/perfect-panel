package user

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateUserNotifyInput struct {
	Body types.UpdateUserNotifyRequest
}

func UpdateUserNotifyHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserNotifyInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserNotifyInput) (*struct{}, error) {
		l := NewUpdateUserNotifyLogic(ctx, svcCtx)
		if err := l.UpdateUserNotify(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserNotifyLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update User Notify
func NewUpdateUserNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserNotifyLogic {
	return &UpdateUserNotifyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserNotifyLogic) UpdateUserNotify(req *types.UpdateUserNotifyRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	if u.Id == 0 {
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "user not login")
	}
	u.EnableLoginNotify = req.EnableLoginNotify
	u.EnableBalanceNotify = req.EnableBalanceNotify
	u.EnableSubscribeNotify = req.EnableSubscribeNotify
	u.EnableTradeNotify = req.EnableTradeNotify
	if err := l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "update user notify error: %v", err.Error())
	}
	return nil
}
