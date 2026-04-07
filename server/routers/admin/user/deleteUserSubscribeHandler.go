// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteUserSubscribeInput struct {
	Body types.DeleteUserSubscribeRequest
}

func DeleteUserSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserSubscribeInput) (*struct{}, error) {
		l := user.NewDeleteUserSubscribeLogic(ctx, svcCtx)
		if err := l.DeleteUserSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
