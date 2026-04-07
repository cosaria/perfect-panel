// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type ResetSortWithServerInput struct {
	Body types.ResetSortRequest
}

func ResetSortWithServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetSortWithServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetSortWithServerInput) (*struct{}, error) {
		l := NewResetSortWithServerLogic(ctx, svcCtx)
		if err := l.ResetSortWithServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
