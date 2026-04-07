package user

import (
	"context"
	"os"
	"strings"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logger.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.GetDetailRequest) error {
	isDemo := strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo"

	if req.Id == 2 && isDemo {
		return errors.Wrapf(xerr.NewErrCodeMsg(503, "Demo mode does not allow deletion of the admin user"), "delete user failed: cannot delete admin user in demo mode")
	}
	err := l.svcCtx.UserModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete user error: %v", err.Error())
	}
	return nil
}
