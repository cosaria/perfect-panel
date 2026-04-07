package user

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetDeviceListOutput struct {
	Body *types.GetDeviceListResponse
}

func GetDeviceListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetDeviceListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetDeviceListOutput, error) {
		l := NewGetDeviceListLogic(ctx, svcCtx)
		resp, err := l.GetDeviceList()
		if err != nil {
			return nil, err
		}
		return &GetDeviceListOutput{Body: resp}, nil
	}
}

type GetDeviceListLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Device List
func NewGetDeviceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceListLogic {
	return &GetDeviceListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceListLogic) GetDeviceList() (resp *types.GetDeviceListResponse, err error) {
	userInfo := l.ctx.Value(config.CtxKeyUser).(*user.User)
	list, count, err := l.svcCtx.UserModel.QueryDeviceList(l.ctx, userInfo.Id)
	userRespList := make([]types.UserDevice, 0)
	tool.DeepCopy(&userRespList, list)
	resp = &types.GetDeviceListResponse{
		Total: count,
		List:  userRespList,
	}
	return
}
