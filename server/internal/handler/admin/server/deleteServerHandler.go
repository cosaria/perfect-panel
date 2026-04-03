// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type DeleteServerInput struct {
	Body types.DeleteServerRequest
}

func DeleteServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteServerInput) (*struct{}, error) {
		l := server.NewDeleteServerLogic(ctx, svcCtx)
		if err := l.DeleteServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
