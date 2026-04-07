package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetOAuthMethodsOutput struct {
	Body *types.GetOAuthMethodsResponse
}

func GetOAuthMethodsHandler(deps Deps) func(context.Context, *struct{}) (*GetOAuthMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetOAuthMethodsOutput, error) {
		l := NewGetOAuthMethodsLogic(ctx, deps)
		resp, err := l.GetOAuthMethods()
		if err != nil {
			return nil, err
		}
		return &GetOAuthMethodsOutput{Body: resp}, nil
	}
}

type GetOAuthMethodsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get OAuth Methods
func NewGetOAuthMethodsLogic(ctx context.Context, deps Deps) *GetOAuthMethodsLogic {
	return &GetOAuthMethodsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetOAuthMethodsLogic) GetOAuthMethods() (resp *types.GetOAuthMethodsResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	methods, err := l.deps.UserModel.FindUserAuthMethods(l.ctx, u.Id)
	if err != nil {
		l.Errorw("find user auth methods failed:", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find user auth methods failed: %v", err.Error())
	}
	list := make([]types.UserAuthMethod, 0)
	tool.DeepCopy(&list, methods)
	return &types.GetOAuthMethodsResponse{
		Methods: list,
	}, nil
}
