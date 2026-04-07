// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UnsubscribeInput struct {
	Body types.UnsubscribeRequest
}

func UnsubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *UnsubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnsubscribeInput) (*struct{}, error) {
		l := user.NewUnsubscribeLogic(ctx, svcCtx)
		if err := l.Unsubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
