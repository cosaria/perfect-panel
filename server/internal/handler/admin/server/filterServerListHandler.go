// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterServerListInput struct {
	types.FilterServerListRequest
}

type FilterServerListOutput struct {
	Body *types.FilterServerListResponse
}

func FilterServerListHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterServerListInput) (*FilterServerListOutput, error) {
	return func(ctx context.Context, input *FilterServerListInput) (*FilterServerListOutput, error) {
		l := server.NewFilterServerListLogic(ctx, svcCtx)
		resp, err := l.FilterServerList(&input.FilterServerListRequest)
		if err != nil {
			return nil, err
		}
		return &FilterServerListOutput{Body: resp}, nil
	}
}
