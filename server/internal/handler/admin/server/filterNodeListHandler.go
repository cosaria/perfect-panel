// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/server"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterNodeListInput struct {
	types.FilterNodeListRequest
}

type FilterNodeListOutput struct {
	Body *types.FilterNodeListResponse
}

func FilterNodeListHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterNodeListInput) (*FilterNodeListOutput, error) {
	return func(ctx context.Context, input *FilterNodeListInput) (*FilterNodeListOutput, error) {
		l := server.NewFilterNodeListLogic(ctx, svcCtx)
		resp, err := l.FilterNodeList(&input.FilterNodeListRequest)
		if err != nil {
			return nil, err
		}
		return &FilterNodeListOutput{Body: resp}, nil
	}
}
