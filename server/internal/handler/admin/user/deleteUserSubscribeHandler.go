// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
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
