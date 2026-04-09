package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type UpdateUserPasswordInput struct {
	Body types.UpdateUserPasswordRequest
}

func UpdateUserPasswordHandler(deps Deps) func(context.Context, *UpdateUserPasswordInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserPasswordInput) (*struct{}, error) {
		l := NewUpdateUserPasswordLogic(ctx, deps)
		if err := l.UpdateUserPassword(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserPasswordLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update User Password
func NewUpdateUserPasswordLogic(ctx context.Context, deps Deps) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdateUserPasswordRequest) error {
	userInfo := l.ctx.Value(config.CtxKeyUser).(*user.User)
	//update the password
	userInfo.Password = tool.EncodePassWord(req.Password)
	if err := l.deps.UserModel.Update(l.ctx, userInfo); err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update user password error")
	}
	return nil
}
