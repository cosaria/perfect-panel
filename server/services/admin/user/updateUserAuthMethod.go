package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateUserAuthMethodInput struct {
	Body types.UpdateUserAuthMethodRequest
}

func UpdateUserAuthMethodHandler(deps Deps) func(context.Context, *UpdateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserAuthMethodInput) (*struct{}, error) {
		l := NewUpdateUserAuthMethodLogic(ctx, deps)
		if err := l.UpdateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserAuthMethodLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update user auth method
func NewUpdateUserAuthMethodLogic(ctx context.Context, deps Deps) *UpdateUserAuthMethodLogic {
	return &UpdateUserAuthMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserAuthMethodLogic) UpdateUserAuthMethod(req *types.UpdateUserAuthMethodRequest) error {
	method, err := l.deps.UserModel.FindUserAuthMethodByPlatform(l.ctx, req.UserId, req.AuthType)
	if err != nil {
		l.Errorw("Get user auth method error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId), logger.Field("authType", req.AuthType))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get user auth method error: %v", err.Error())
	}
	userInfo, err := l.deps.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("Get user info error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get user info error: %v", err.Error())
	}

	method.AuthType = req.AuthType
	method.AuthIdentifier = req.AuthIdentifier
	if err = l.deps.UserModel.UpdateUserAuthMethods(l.ctx, method); err != nil {
		l.Errorw("Update user auth method error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId), logger.Field("authType", req.AuthType))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update user auth method error: %v", err.Error())
	}
	if err = l.deps.UserModel.UpdateUserCache(l.ctx, userInfo); err != nil {
		l.Errorw("Update user cache error", logger.Field("error", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Update user cache error: %v", err.Error())
	}
	return nil
}
