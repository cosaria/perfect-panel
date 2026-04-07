package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserSubscribeDevicesInput struct {
	types.GetUserSubscribeDevicesRequest
}

type GetUserSubscribeDevicesOutput struct {
	Body *types.GetUserSubscribeDevicesResponse
}

func GetUserSubscribeDevicesHandler(deps Deps) func(context.Context, *GetUserSubscribeDevicesInput) (*GetUserSubscribeDevicesOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeDevicesInput) (*GetUserSubscribeDevicesOutput, error) {
		l := NewGetUserSubscribeDevicesLogic(ctx, deps)
		resp, err := l.GetUserSubscribeDevices(&input.GetUserSubscribeDevicesRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeDevicesOutput{Body: resp}, nil
	}
}

type GetUserSubscribeDevicesLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe devices
func NewGetUserSubscribeDevicesLogic(ctx context.Context, deps Deps) *GetUserSubscribeDevicesLogic {
	return &GetUserSubscribeDevicesLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeDevicesLogic) GetUserSubscribeDevices(req *types.GetUserSubscribeDevicesRequest) (resp *types.GetUserSubscribeDevicesResponse, err error) {
	list, total, err := l.deps.UserModel.QueryDevicePageList(l.ctx, req.UserId, req.SubscribeId, req.Page, req.Size)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetUserSubscribeDevices failed: %v", err.Error())
	}
	userRespList := make([]types.UserDevice, 0)
	tool.DeepCopy(&userRespList, list)
	return &types.GetUserSubscribeDevicesResponse{
		Total: total,
		List:  userRespList,
	}, nil
}
