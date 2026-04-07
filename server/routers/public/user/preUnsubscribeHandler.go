// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PreUnsubscribeInput struct {
	Body types.PreUnsubscribeRequest
}

type PreUnsubscribeOutput struct {
	Body *types.PreUnsubscribeResponse
}

func PreUnsubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *PreUnsubscribeInput) (*PreUnsubscribeOutput, error) {
	return func(ctx context.Context, input *PreUnsubscribeInput) (*PreUnsubscribeOutput, error) {
		l := user.NewPreUnsubscribeLogic(ctx, svcCtx)
		resp, err := l.PreUnsubscribe(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PreUnsubscribeOutput{Body: resp}, nil
	}
}
