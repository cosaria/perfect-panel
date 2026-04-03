// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/document"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type BatchDeleteDocumentInput struct {
	Body types.BatchDeleteDocumentRequest
}

func BatchDeleteDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteDocumentInput) (*struct{}, error) {
		l := document.NewBatchDeleteDocumentLogic(ctx, svcCtx)
		if err := l.BatchDeleteDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
