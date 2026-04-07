// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserSubscribeInput struct {
	Body types.UpdateUserSubscribeRequest
}

func UpdateUserSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserSubscribeInput) (*struct{}, error) {
		l := NewUpdateUserSubscribeLogic(ctx, svcCtx)
		if err := l.UpdateUserSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
