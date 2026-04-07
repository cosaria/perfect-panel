// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetMessageLogListInput struct {
	types.GetMessageLogListRequest
}

type GetMessageLogListOutput struct {
	Body *types.GetMessageLogListResponse
}

func GetMessageLogListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetMessageLogListInput) (*GetMessageLogListOutput, error) {
	return func(ctx context.Context, input *GetMessageLogListInput) (*GetMessageLogListOutput, error) {
		l := NewGetMessageLogListLogic(ctx, svcCtx)
		resp, err := l.GetMessageLogList(&input.GetMessageLogListRequest)
		if err != nil {
			return nil, err
		}
		return &GetMessageLogListOutput{Body: resp}, nil
	}
}
