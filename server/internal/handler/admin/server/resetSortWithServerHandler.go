// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type ResetSortWithServerInput struct {
	Body types.ResetSortRequest
}

func ResetSortWithServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetSortWithServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetSortWithServerInput) (*struct{}, error) {
		l := server.NewResetSortWithServerLogic(ctx, svcCtx)
		if err := l.ResetSortWithServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
