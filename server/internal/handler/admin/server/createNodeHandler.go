// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateNodeInput struct {
	Body types.CreateNodeRequest
}

func CreateNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateNodeInput) (*struct{}, error) {
		l := server.NewCreateNodeLogic(ctx, svcCtx)
		if err := l.CreateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
