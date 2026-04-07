// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type ToggleUserSubscribeStatusInput struct {
	Body types.ToggleUserSubscribeStatusRequest
}

func ToggleUserSubscribeStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *ToggleUserSubscribeStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *ToggleUserSubscribeStatusInput) (*struct{}, error) {
		l := NewToggleUserSubscribeStatusLogic(ctx, svcCtx)
		if err := l.ToggleUserSubscribeStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
