package user

import (
	"context"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserPasswordLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update User Password
func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdateUserPasswordRequest) error {
	userInfo := l.ctx.Value(config.CtxKeyUser).(*user.User)
	//update the password
	userInfo.Password = tool.EncodePassWord(req.Password)
	if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update user password error")
	}
	return nil
}
