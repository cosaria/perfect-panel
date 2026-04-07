// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetClientOutput struct {
	Body *types.GetSubscribeClientResponse
}

func GetClientHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetClientOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetClientOutput, error) {
		l := NewGetClientLogic(ctx, svcCtx)
		resp, err := l.GetClient()
		if err != nil {
			return nil, err
		}
		return &GetClientOutput{Body: resp}, nil
	}
}
