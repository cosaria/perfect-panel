package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
)

type GetDeviceListOutput struct {
	Body *types.GetDeviceListResponse
}

func GetDeviceListHandler(deps Deps) func(context.Context, *struct{}) (*GetDeviceListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetDeviceListOutput, error) {
		l := NewGetDeviceListLogic(ctx, deps)
		resp, err := l.GetDeviceList()
		if err != nil {
			return nil, err
		}
		return &GetDeviceListOutput{Body: resp}, nil
	}
}

type GetDeviceListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Device List
func NewGetDeviceListLogic(ctx context.Context, deps Deps) *GetDeviceListLogic {
	return &GetDeviceListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetDeviceListLogic) GetDeviceList() (resp *types.GetDeviceListResponse, err error) {
	userInfo := l.ctx.Value(config.CtxKeyUser).(*user.User)
	list, count, err := l.deps.UserModel.QueryDeviceList(l.ctx, userInfo.Id)
	userRespList := make([]types.UserDevice, 0)
	tool.DeepCopy(&userRespList, list)
	resp = &types.GetDeviceListResponse{
		Total: count,
		List:  userRespList,
	}
	return
}
