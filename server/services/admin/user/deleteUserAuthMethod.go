package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteUserAuthMethodInput struct {
	Body types.DeleteUserAuthMethodRequest
}

func DeleteUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserAuthMethodInput) (*struct{}, error) {
		l := NewDeleteUserAuthMethodLogic(ctx, svcCtx)
		if err := l.DeleteUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteUserAuthMethodLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete user auth method
func NewDeleteUserAuthMethodLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserAuthMethodLogic {
	return &DeleteUserAuthMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserAuthMethodLogic) DeleteUserAuthMethod(req *types.DeleteUserAuthMethodRequest) error {
	err := l.svcCtx.UserModel.DeleteUserAuthMethods(l.ctx, req.UserId, req.AuthType)
	if err != nil {
		l.Errorw("[DeleteUserAuthMethodLogic] Delete User Auth Method Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId), logger.Field("authType", req.AuthType))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "Delete User Auth Method Error")
	}
	return nil
}
