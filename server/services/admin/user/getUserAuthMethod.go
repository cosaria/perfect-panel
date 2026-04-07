package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserAuthMethodInput struct {
	types.GetUserAuthMethodRequest
}

type GetUserAuthMethodOutput struct {
	Body *types.GetUserAuthMethodResponse
}

func GetUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserAuthMethodInput) (*GetUserAuthMethodOutput, error) {
	return func(ctx context.Context, input *GetUserAuthMethodInput) (*GetUserAuthMethodOutput, error) {
		l := NewGetUserAuthMethodLogic(ctx, svcCtx)
		resp, err := l.GetUserAuthMethod(&input.GetUserAuthMethodRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserAuthMethodOutput{Body: resp}, nil
	}
}

type GetUserAuthMethodLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get user auth method
func NewGetUserAuthMethodLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthMethodLogic {
	return &GetUserAuthMethodLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAuthMethodLogic) GetUserAuthMethod(req *types.GetUserAuthMethodRequest) (resp *types.GetUserAuthMethodResponse, err error) {
	methods, err := l.svcCtx.UserModel.FindUserAuthMethods(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("[GetUserAuthMethodLogic] Get User Auth Method Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get User Auth Method Error")
	}
	list := make([]types.UserAuthMethod, 0)
	tool.DeepCopy(&list, methods)

	return &types.GetUserAuthMethodResponse{
		AuthMethods: list,
	}, nil
}
