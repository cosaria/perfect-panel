// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSubscriptionInput struct {
	types.GetSubscriptionRequest
}

type GetSubscriptionOutput struct {
	Body *types.GetSubscriptionResponse
}

func GetSubscriptionHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscriptionInput) (*GetSubscriptionOutput, error) {
	return func(ctx context.Context, input *GetSubscriptionInput) (*GetSubscriptionOutput, error) {
		l := NewGetSubscriptionLogic(ctx, svcCtx)
		resp, err := l.GetSubscription(&input.GetSubscriptionRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscriptionOutput{Body: resp}, nil
	}
}
