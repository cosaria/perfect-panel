// huma:migrated
package user

import (
	"context"
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
