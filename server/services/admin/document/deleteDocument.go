package document

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteDocumentInput struct {
	Body types.DeleteDocumentRequest
}

func DeleteDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteDocumentInput) (*struct{}, error) {
		l := NewDeleteDocumentLogic(ctx, svcCtx)
		if err := l.DeleteDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteDocumentLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete document
func NewDeleteDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocumentLogic {
	return &DeleteDocumentLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDocumentLogic) DeleteDocument(req *types.DeleteDocumentRequest) error {
	if err := l.svcCtx.DocumentModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("[DeleteDocument] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "failed to delete document: %v", err.Error())
	}
	return nil
}
