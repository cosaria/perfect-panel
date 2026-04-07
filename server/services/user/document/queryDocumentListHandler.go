// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryDocumentListOutput struct {
	Body *types.QueryDocumentListResponse
}

func QueryDocumentListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryDocumentListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryDocumentListOutput, error) {
		l := NewQueryDocumentListLogic(ctx, svcCtx)
		resp, err := l.QueryDocumentList()
		if err != nil {
			return nil, err
		}
		return &QueryDocumentListOutput{Body: resp}, nil
	}
}
