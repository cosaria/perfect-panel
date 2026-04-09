package auth

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CheckUserInput struct {
	types.CheckUserRequest
}

type CheckUserOutput struct {
	Body *types.CheckUserResponse
}

func CheckUserHandler(deps Deps) func(context.Context, *CheckUserInput) (*CheckUserOutput, error) {
	return func(ctx context.Context, input *CheckUserInput) (*CheckUserOutput, error) {
		l := NewCheckUserLogic(ctx, deps)
		resp, err := l.CheckUser(&input.CheckUserRequest)
		if err != nil {
			return nil, err
		}
		return &CheckUserOutput{Body: resp}, nil
	}
}

type CheckUserLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewCheckUserLogic Check user is exist
func NewCheckUserLogic(ctx context.Context, deps Deps) *CheckUserLogic {
	return &CheckUserLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CheckUserLogic) CheckUser(req *types.CheckUserRequest) (resp *types.CheckUserResponse, err error) {
	authMethod, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, "email", req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user by email error: %v", err.Error())
	}
	return &types.CheckUserResponse{
		Exist: authMethod.UserId != 0,
	}, nil
}
