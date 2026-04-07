// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/services/common"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetGlobalConfigOutput struct {
	Body *types.GetGlobalConfigResponse
}

func GetGlobalConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetGlobalConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetGlobalConfigOutput, error) {
		l := common.NewGetGlobalConfigLogic(ctx, svcCtx)
		resp, err := l.GetGlobalConfig()
		if err != nil {
			return nil, err
		}
		return &GetGlobalConfigOutput{Body: resp}, nil
	}
}
