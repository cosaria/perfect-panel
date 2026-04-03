// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserSubscribeDevicesInput struct {
	types.GetUserSubscribeDevicesRequest
}

type GetUserSubscribeDevicesOutput struct {
	Body *types.GetUserSubscribeDevicesResponse
}

func GetUserSubscribeDevicesHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeDevicesInput) (*GetUserSubscribeDevicesOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeDevicesInput) (*GetUserSubscribeDevicesOutput, error) {
		l := user.NewGetUserSubscribeDevicesLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeDevices(&input.GetUserSubscribeDevicesRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeDevicesOutput{Body: resp}, nil
	}
}
