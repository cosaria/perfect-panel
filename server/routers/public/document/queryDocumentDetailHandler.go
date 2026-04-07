// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/services/user/document"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryDocumentDetailInput struct {
	types.QueryDocumentDetailRequest
}

type QueryDocumentDetailOutput struct {
	Body *types.Document
}

func QueryDocumentDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryDocumentDetailInput) (*QueryDocumentDetailOutput, error) {
	return func(ctx context.Context, input *QueryDocumentDetailInput) (*QueryDocumentDetailOutput, error) {
		l := document.NewQueryDocumentDetailLogic(ctx, svcCtx)
		resp, err := l.QueryDocumentDetail(&input.QueryDocumentDetailRequest)
		if err != nil {
			return nil, err
		}
		return &QueryDocumentDetailOutput{Body: resp}, nil
	}
}
