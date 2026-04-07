// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSubscribeDetailsInput struct {
	types.GetSubscribeDetailsRequest
}

type GetSubscribeDetailsOutput struct {
	Body *types.Subscribe
}

func GetSubscribeDetailsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscribeDetailsInput) (*GetSubscribeDetailsOutput, error) {
	return func(ctx context.Context, input *GetSubscribeDetailsInput) (*GetSubscribeDetailsOutput, error) {
		l := NewGetSubscribeDetailsLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeDetails(&input.GetSubscribeDetailsRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeDetailsOutput{Body: resp}, nil
	}
}
