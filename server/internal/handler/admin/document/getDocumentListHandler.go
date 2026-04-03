// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/document"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetDocumentListInput struct {
	types.GetDocumentListRequest
}

type GetDocumentListOutput struct {
	Body *types.GetDocumentListResponse
}

func GetDocumentListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetDocumentListInput) (*GetDocumentListOutput, error) {
	return func(ctx context.Context, input *GetDocumentListInput) (*GetDocumentListOutput, error) {
		l := document.NewGetDocumentListLogic(ctx, svcCtx)
		resp, err := l.GetDocumentList(&input.GetDocumentListRequest)
		if err != nil {
			return nil, err
		}
		return &GetDocumentListOutput{Body: resp}, nil
	}
}
