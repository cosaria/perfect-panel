package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type DeleteUserAuthMethodInput struct {
	Body types.DeleteUserAuthMethodRequest
}

func DeleteUserAuthMethodHandler(deps Deps) func(context.Context, *DeleteUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserAuthMethodInput) (*struct{}, error) {
		l := NewDeleteUserAuthMethodLogic(ctx, deps)
		if err := l.DeleteUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteUserAuthMethodLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete user auth method
func NewDeleteUserAuthMethodLogic(ctx context.Context, deps Deps) *DeleteUserAuthMethodLogic {
	return &DeleteUserAuthMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteUserAuthMethodLogic) DeleteUserAuthMethod(req *types.DeleteUserAuthMethodRequest) error {
	err := l.deps.UserModel.DeleteUserAuthMethods(l.ctx, req.UserId, req.AuthType)
	if err != nil {
		l.Errorw("[DeleteUserAuthMethodLogic] Delete User Auth Method Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId), logger.Field("authType", req.AuthType))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "Delete User Auth Method Error")
	}
	return nil
}
