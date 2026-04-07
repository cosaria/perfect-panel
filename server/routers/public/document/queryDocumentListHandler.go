// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/services/user/document"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryDocumentListOutput struct {
	Body *types.QueryDocumentListResponse
}

func QueryDocumentListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryDocumentListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryDocumentListOutput, error) {
		l := document.NewQueryDocumentListLogic(ctx, svcCtx)
		resp, err := l.QueryDocumentList()
		if err != nil {
			return nil, err
		}
		return &QueryDocumentListOutput{Body: resp}, nil
	}
}
