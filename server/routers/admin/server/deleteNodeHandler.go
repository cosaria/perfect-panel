// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/services/admin/server"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteNodeInput struct {
	Body types.DeleteNodeRequest
}

func DeleteNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteNodeInput) (*struct{}, error) {
		l := server.NewDeleteNodeLogic(ctx, svcCtx)
		if err := l.DeleteNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
