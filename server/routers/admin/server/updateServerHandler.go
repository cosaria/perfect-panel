// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/services/admin/server"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateServerInput struct {
	Body types.UpdateServerRequest
}

func UpdateServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateServerInput) (*struct{}, error) {
		l := server.NewUpdateServerLogic(ctx, svcCtx)
		if err := l.UpdateServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
