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

type GetUserDetailInput struct {
	types.GetDetailRequest
}

type GetUserDetailOutput struct {
	Body *types.User
}

func GetUserDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserDetailInput) (*GetUserDetailOutput, error) {
	return func(ctx context.Context, input *GetUserDetailInput) (*GetUserDetailOutput, error) {
		l := NewGetUserDetailLogic(ctx, svcCtx)
		resp, err := l.GetUserDetail(&input.GetDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserDetailOutput{Body: resp}, nil
	}
}

type GetUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger.Logger
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logger.WithContext(ctx),
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.GetDetailRequest) (*types.User, error) {
	resp := types.User{}
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get user detail error: %v", err.Error())
	}
	tool.DeepCopy(&resp, userInfo)
	return &resp, nil
}
