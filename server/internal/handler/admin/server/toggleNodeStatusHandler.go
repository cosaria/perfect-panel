// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type ToggleNodeStatusInput struct {
	Body types.ToggleNodeStatusRequest
}

func ToggleNodeStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *ToggleNodeStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *ToggleNodeStatusInput) (*struct{}, error) {
		l := server.NewToggleNodeStatusLogic(ctx, svcCtx)
		if err := l.ToggleNodeStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
