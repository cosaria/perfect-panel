package document

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type BatchDeleteDocumentInput struct {
	Body types.BatchDeleteDocumentRequest
}

func BatchDeleteDocumentHandler(deps Deps) func(context.Context, *BatchDeleteDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteDocumentInput) (*struct{}, error) {
		l := NewBatchDeleteDocumentLogic(ctx, deps)
		if err := l.BatchDeleteDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteDocumentLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Batch delete document
func NewBatchDeleteDocumentLogic(ctx context.Context, deps Deps) *BatchDeleteDocumentLogic {
	return &BatchDeleteDocumentLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *BatchDeleteDocumentLogic) BatchDeleteDocument(req *types.BatchDeleteDocumentRequest) error {
	for _, id := range req.Ids {
		if err := l.deps.DocumentModel.Delete(l.ctx, id); err != nil {
			l.Errorw("[BatchDeleteDocument] Database Error", logger.Field("error", err.Error()))
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "failed to delete document: %v", err.Error())
		}
	}
	return nil
}
