// huma:migrated
package document

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateDocumentInput struct {
	Body types.UpdateDocumentRequest
}

func UpdateDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateDocumentInput) (*struct{}, error) {
		l := NewUpdateDocumentLogic(ctx, svcCtx)
		if err := l.UpdateDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
