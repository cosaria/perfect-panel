// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/services/admin/document"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteDocumentInput struct {
	Body types.DeleteDocumentRequest
}

func DeleteDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteDocumentInput) (*struct{}, error) {
		l := document.NewDeleteDocumentLogic(ctx, svcCtx)
		if err := l.DeleteDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
