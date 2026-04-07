// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSubscribeConfigOutput struct {
	Body *types.SubscribeConfig
}

func GetSubscribeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSubscribeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSubscribeConfigOutput, error) {
		l := NewGetSubscribeConfigLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeConfig()
		if err != nil {
			return nil, err
		}
		return &GetSubscribeConfigOutput{Body: resp}, nil
	}
}
