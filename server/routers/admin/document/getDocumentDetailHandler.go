// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/services/admin/document"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetDocumentDetailInput struct {
	types.GetDocumentDetailRequest
}

type GetDocumentDetailOutput struct {
	Body *types.Document
}

func GetDocumentDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetDocumentDetailInput) (*GetDocumentDetailOutput, error) {
	return func(ctx context.Context, input *GetDocumentDetailInput) (*GetDocumentDetailOutput, error) {
		l := document.NewGetDocumentDetailLogic(ctx, svcCtx)
		resp, err := l.GetDocumentDetail(&input.GetDocumentDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetDocumentDetailOutput{Body: resp}, nil
	}
}
