// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/services/admin/server"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetServerProtocolsInput struct {
	types.GetServerProtocolsRequest
}

type GetServerProtocolsOutput struct {
	Body *types.GetServerProtocolsResponse
}

func GetServerProtocolsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
	return func(ctx context.Context, input *GetServerProtocolsInput) (*GetServerProtocolsOutput, error) {
		l := server.NewGetServerProtocolsLogic(ctx, svcCtx)
		resp, err := l.GetServerProtocols(&input.GetServerProtocolsRequest)
		if err != nil {
			return nil, err
		}
		return &GetServerProtocolsOutput{Body: resp}, nil
	}
}
