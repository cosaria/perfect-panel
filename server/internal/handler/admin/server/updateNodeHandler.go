// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateNodeInput struct {
	Body types.UpdateNodeRequest
}

func UpdateNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateNodeInput) (*struct{}, error) {
		l := server.NewUpdateNodeLogic(ctx, svcCtx)
		if err := l.UpdateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
