// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type BatchDeleteUserInput struct {
	Body types.BatchDeleteUserRequest
}

func BatchDeleteUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteUserInput) (*struct{}, error) {
		l := user.NewBatchDeleteUserLogic(ctx, svcCtx)
		if err := l.BatchDeleteUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
