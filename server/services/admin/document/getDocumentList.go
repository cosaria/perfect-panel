package document

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetDocumentListInput struct {
	types.GetDocumentListRequest
}

type GetDocumentListOutput struct {
	Body *types.GetDocumentListResponse
}

func GetDocumentListHandler(deps Deps) func(context.Context, *GetDocumentListInput) (*GetDocumentListOutput, error) {
	return func(ctx context.Context, input *GetDocumentListInput) (*GetDocumentListOutput, error) {
		l := NewGetDocumentListLogic(ctx, deps)
		resp, err := l.GetDocumentList(&input.GetDocumentListRequest)
		if err != nil {
			return nil, err
		}
		return &GetDocumentListOutput{Body: resp}, nil
	}
}

type GetDocumentListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get document list
func NewGetDocumentListLogic(ctx context.Context, deps Deps) *GetDocumentListLogic {
	return &GetDocumentListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetDocumentListLogic) GetDocumentList(req *types.GetDocumentListRequest) (resp *types.GetDocumentListResponse, err error) {
	total, data, err := l.deps.DocumentModel.QueryDocumentList(l.ctx, int(req.Page), int(req.Size), req.Tag, req.Search)
	if err != nil {
		l.Errorw("[GetDocumentList] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryDocumentList error: %v", err.Error())
	}
	resp = &types.GetDocumentListResponse{
		Total: total,
		List:  make([]types.Document, 0),
	}
	for _, v := range data {
		resp.List = append(resp.List, types.Document{
			Id:        v.Id,
			Title:     v.Title,
			Tags:      tool.StringMergeAndRemoveDuplicates(v.Tags),
			Content:   v.Content,
			Show:      *v.Show,
			CreatedAt: v.CreatedAt.UnixMilli(),
			UpdatedAt: v.UpdatedAt.UnixMilli(),
		})
	}
	return
}
