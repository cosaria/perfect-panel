package server

import (
	"context"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"strings"
)

type FilterNodeListInput struct {
	types.FilterNodeListRequest
}

type FilterNodeListOutput struct {
	Body *types.FilterNodeListResponse
}

func FilterNodeListHandler(deps Deps) func(context.Context, *FilterNodeListInput) (*FilterNodeListOutput, error) {
	return func(ctx context.Context, input *FilterNodeListInput) (*FilterNodeListOutput, error) {
		l := NewFilterNodeListLogic(ctx, deps)
		resp, err := l.FilterNodeList(&input.FilterNodeListRequest)
		if err != nil {
			return nil, err
		}
		return &FilterNodeListOutput{Body: resp}, nil
	}
}

type FilterNodeListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewFilterNodeListLogic Filter Node List
func NewFilterNodeListLogic(ctx context.Context, deps Deps) *FilterNodeListLogic {
	return &FilterNodeListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *FilterNodeListLogic) FilterNodeList(req *types.FilterNodeListRequest) (resp *types.FilterNodeListResponse, err error) {
	total, data, err := l.deps.NodeModel.FilterNodeList(l.ctx, &node.FilterNodeParams{
		Page:   req.Page,
		Size:   req.Size,
		Search: req.Search,
	})

	if err != nil {
		l.Errorw("[FilterNodeList] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[FilterNodeList] Query Database Error")
	}

	list := make([]types.Node, 0)
	for _, datum := range data {
		list = append(list, types.Node{
			Id:        datum.Id,
			Name:      datum.Name,
			Tags:      tool.RemoveDuplicateElements(strings.Split(datum.Tags, ",")...),
			Port:      datum.Port,
			Address:   datum.Address,
			ServerId:  datum.ServerId,
			Protocol:  datum.Protocol,
			Enabled:   datum.Enabled,
			Sort:      datum.Sort,
			CreatedAt: datum.CreatedAt.UnixMilli(),
			UpdatedAt: datum.UpdatedAt.UnixMilli(),
		})
	}

	return &types.FilterNodeListResponse{
		List:  list,
		Total: total,
	}, nil
}
