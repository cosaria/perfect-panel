// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/document"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateDocumentInput struct {
	Body types.CreateDocumentRequest
}

func CreateDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateDocumentInput) (*struct{}, error) {
		l := document.NewCreateDocumentLogic(ctx, svcCtx)
		if err := l.CreateDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
