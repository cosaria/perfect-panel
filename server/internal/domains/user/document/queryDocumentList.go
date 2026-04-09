package document

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type QueryDocumentListOutput struct {
	Body *types.QueryDocumentListResponse
}

func QueryDocumentListHandler(deps Deps) func(context.Context, *struct{}) (*QueryDocumentListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryDocumentListOutput, error) {
		l := NewQueryDocumentListLogic(ctx, deps)
		resp, err := l.QueryDocumentList()
		if err != nil {
			return nil, err
		}
		return &QueryDocumentListOutput{Body: resp}, nil
	}
}

type QueryDocumentListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get document list
func NewQueryDocumentListLogic(ctx context.Context, deps Deps) *QueryDocumentListLogic {
	return &QueryDocumentListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryDocumentListLogic) QueryDocumentList() (resp *types.QueryDocumentListResponse, err error) {
	total, data, err := l.deps.DocumentModel.GetDocumentListByAll(l.ctx)
	if err != nil {
		l.Errorw("[QueryDocumentList] error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryDocumentList error: %v", err.Error())
	}
	resp = &types.QueryDocumentListResponse{
		Total: total,
		List:  make([]types.Document, 0),
	}
	for _, item := range data {
		resp.List = append(resp.List, types.Document{
			Id:        item.Id,
			Title:     item.Title,
			Tags:      tool.StringMergeAndRemoveDuplicates(item.Tags),
			UpdatedAt: item.UpdatedAt.UnixMilli(),
		})
	}
	return
}
