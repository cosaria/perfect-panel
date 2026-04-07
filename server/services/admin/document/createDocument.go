package document

import (
	"context"
	"github.com/perfect-panel/server/models/document"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"strings"
)

type CreateDocumentInput struct {
	Body types.CreateDocumentRequest
}

func CreateDocumentHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateDocumentInput) (*struct{}, error) {
		l := NewCreateDocumentLogic(ctx, svcCtx)
		if err := l.CreateDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateDocumentLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create document
func NewCreateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocumentLogic {
	return &CreateDocumentLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDocumentLogic) CreateDocument(req *types.CreateDocumentRequest) error {
	if err := l.svcCtx.DocumentModel.Insert(l.ctx, &document.Document{
		Title:   req.Title,
		Content: req.Content,
		Tags:    strings.Join(req.Tags, ","),
		Show:    req.Show,
	}); err != nil {
		l.Errorw("[CreateDocument] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert document error: %v", err.Error())
	}
	return nil
}
