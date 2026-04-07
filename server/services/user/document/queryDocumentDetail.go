package document

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type QueryDocumentDetailInput struct {
	types.QueryDocumentDetailRequest
}

type QueryDocumentDetailOutput struct {
	Body *types.Document
}

func QueryDocumentDetailHandler(deps Deps) func(context.Context, *QueryDocumentDetailInput) (*QueryDocumentDetailOutput, error) {
	return func(ctx context.Context, input *QueryDocumentDetailInput) (*QueryDocumentDetailOutput, error) {
		l := NewQueryDocumentDetailLogic(ctx, deps)
		resp, err := l.QueryDocumentDetail(&input.QueryDocumentDetailRequest)
		if err != nil {
			return nil, err
		}
		return &QueryDocumentDetailOutput{Body: resp}, nil
	}
}

type QueryDocumentDetailLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get document detail
func NewQueryDocumentDetailLogic(ctx context.Context, deps Deps) *QueryDocumentDetailLogic {
	return &QueryDocumentDetailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryDocumentDetailLogic) QueryDocumentDetail(req *types.QueryDocumentDetailRequest) (resp *types.Document, err error) {
	// find document
	data, err := l.deps.DocumentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[QueryDocumentDetailLogic] FindOne error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %s", err.Error())
	}
	resp = &types.Document{}
	tool.DeepCopy(resp, data)
	return
}
