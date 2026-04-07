// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type ResetSortWithNodeInput struct {
	Body types.ResetSortRequest
}

func ResetSortWithNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetSortWithNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetSortWithNodeInput) (*struct{}, error) {
		l := NewResetSortWithNodeLogic(ctx, svcCtx)
		if err := l.ResetSortWithNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
