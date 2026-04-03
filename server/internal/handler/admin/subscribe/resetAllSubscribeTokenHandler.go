// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type ResetAllSubscribeTokenOutput struct {
	Body *types.ResetAllSubscribeTokenResponse
}

func ResetAllSubscribeTokenHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*ResetAllSubscribeTokenOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*ResetAllSubscribeTokenOutput, error) {
		l := subscribe.NewResetAllSubscribeTokenLogic(ctx, svcCtx)
		resp, err := l.ResetAllSubscribeToken()
		if err != nil {
			return nil, err
		}
		return &ResetAllSubscribeTokenOutput{Body: resp}, nil
	}
}
