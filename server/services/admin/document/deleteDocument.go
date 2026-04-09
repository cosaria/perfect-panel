package document

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type DeleteDocumentInput struct {
	Body types.DeleteDocumentRequest
}

func DeleteDocumentHandler(deps Deps) func(context.Context, *DeleteDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteDocumentInput) (*struct{}, error) {
		l := NewDeleteDocumentLogic(ctx, deps)
		if err := l.DeleteDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteDocumentLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete document
func NewDeleteDocumentLogic(ctx context.Context, deps Deps) *DeleteDocumentLogic {
	return &DeleteDocumentLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteDocumentLogic) DeleteDocument(req *types.DeleteDocumentRequest) error {
	if err := l.deps.DocumentModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("[DeleteDocument] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "failed to delete document: %v", err.Error())
	}
	return nil
}
