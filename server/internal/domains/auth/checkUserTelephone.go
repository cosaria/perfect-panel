package auth

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/phone"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CheckUserTelephoneInput struct {
	types.TelephoneCheckUserRequest
}

type CheckUserTelephoneOutput struct {
	Body *types.TelephoneCheckUserResponse
}

func CheckUserTelephoneHandler(deps Deps) func(context.Context, *CheckUserTelephoneInput) (*CheckUserTelephoneOutput, error) {
	return func(ctx context.Context, input *CheckUserTelephoneInput) (*CheckUserTelephoneOutput, error) {
		l := NewCheckUserTelephoneLogic(ctx, deps)
		resp, err := l.CheckUserTelephone(&input.TelephoneCheckUserRequest)
		if err != nil {
			return nil, err
		}
		return &CheckUserTelephoneOutput{Body: resp}, nil
	}
}

type CheckUserTelephoneLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Check user telephone is exist
func NewCheckUserTelephoneLogic(ctx context.Context, deps Deps) *CheckUserTelephoneLogic {
	return &CheckUserTelephoneLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CheckUserTelephoneLogic) CheckUserTelephone(req *types.TelephoneCheckUserRequest) (resp *types.TelephoneCheckUserResponse, err error) {
	phoneNumber, err := phone.FormatToE164(req.TelephoneAreaCode, req.Telephone)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TelephoneError), "Invalid phone number")
	}
	authMethods, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "mobile", phoneNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user by email error: %v", err.Error())
	}

	return &types.TelephoneCheckUserResponse{
		Exist: authMethods.UserId != 0,
	}, nil
}
