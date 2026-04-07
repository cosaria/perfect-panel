package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateBindEmailInput struct {
	Body types.UpdateBindEmailRequest
}

func UpdateBindEmailHandler(deps Deps) func(context.Context, *UpdateBindEmailInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateBindEmailInput) (*struct{}, error) {
		l := NewUpdateBindEmailLogic(ctx, deps)
		if err := l.UpdateBindEmail(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateBindEmailLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateBindEmailLogic Update Bind Email
func NewUpdateBindEmailLogic(ctx context.Context, deps Deps) *UpdateBindEmailLogic {
	return &UpdateBindEmailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateBindEmailLogic) UpdateBindEmail(req *types.UpdateBindEmailRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	method, err := l.deps.UserModel.FindUserAuthMethodByUserId(l.ctx, "email", u.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	m, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindUserAuthMethodByOpenID error")
	}
	// email already bind
	if m.Id > 0 {
		return errors.Wrapf(xerr.NewErrCode(xerr.UserExist), "email already bind")
	}
	if method.Id == 0 {
		method = &user.AuthMethods{
			UserId:         u.Id,
			AuthType:       "email",
			AuthIdentifier: req.Email,
			Verified:       false,
		}
		if err := l.deps.UserModel.InsertUserAuthMethods(l.ctx, method); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "InsertUserAuthMethods error")
		}
	} else {
		method.Verified = false
		method.AuthIdentifier = req.Email
		if err := l.deps.UserModel.UpdateUserAuthMethods(l.ctx, method); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateUserAuthMethods error")
		}
	}
	return nil
}
