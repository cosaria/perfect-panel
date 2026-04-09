package document

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/document"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type UpdateDocumentInput struct {
	Body types.UpdateDocumentRequest
}

func UpdateDocumentHandler(deps Deps) func(context.Context, *UpdateDocumentInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateDocumentInput) (*struct{}, error) {
		l := NewUpdateDocumentLogic(ctx, deps)
		if err := l.UpdateDocument(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateDocumentLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update document
func NewUpdateDocumentLogic(ctx context.Context, deps Deps) *UpdateDocumentLogic {
	return &UpdateDocumentLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateDocumentLogic) UpdateDocument(req *types.UpdateDocumentRequest) error {
	if err := l.deps.DocumentModel.Update(l.ctx, &document.Document{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Tags:    strings.Join(req.Tags, ","),
		Show:    req.Show,
	}); err != nil {
		l.Errorw("[UpdateDocument] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "failed to update document: %v", err.Error())
	}
	return nil
}
