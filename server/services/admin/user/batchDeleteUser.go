package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type BatchDeleteUserInput struct {
	Body types.BatchDeleteUserRequest
}

func BatchDeleteUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteUserInput) (*struct{}, error) {
		l := NewBatchDeleteUserLogic(ctx, svcCtx)
		if err := l.BatchDeleteUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger.Logger
}

func NewBatchDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteUserLogic {
	return &BatchDeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logger.WithContext(ctx),
	}
}

func (l *BatchDeleteUserLogic) BatchDeleteUser(req *types.BatchDeleteUserRequest) error {
	isDemo := strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo"

	if tool.Contains(req.Ids, 2) && isDemo {
		return errors.Wrapf(xerr.NewErrCodeMsg(503, "Demo mode does not allow deletion of the admin user"), "BatchDeleteUser failed: cannot delete admin user in demo mode")
	}

	err := l.svcCtx.UserModel.BatchDeleteUser(l.ctx, req.Ids)
	if err != nil {
		l.Logger.Error("[BatchDeleteUserLogic] BatchDeleteUser failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "BatchDeleteUser failed: %v", err.Error())
	}
	return nil
}
