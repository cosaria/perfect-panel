package order

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetOrderListInput struct {
	types.GetOrderListRequest
}

type GetOrderListOutput struct {
	Body *types.GetOrderListResponse
}

func GetOrderListHandler(deps Deps) func(context.Context, *GetOrderListInput) (*GetOrderListOutput, error) {
	return func(ctx context.Context, input *GetOrderListInput) (*GetOrderListOutput, error) {
		l := NewGetOrderListLogic(ctx, deps)
		resp, err := l.GetOrderList(&input.GetOrderListRequest)
		if err != nil {
			return nil, err
		}
		return &GetOrderListOutput{Body: resp}, nil
	}
}

type GetOrderListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetOrderListLogic Get order list
func NewGetOrderListLogic(ctx context.Context, deps Deps) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListRequest) (resp *types.GetOrderListResponse, err error) {
	total, list, err := l.deps.OrderModel.QueryOrderListByPage(l.ctx, int(req.Page), int(req.Size), req.Status, req.UserId, req.SubscribeId, req.Search)
	if err != nil {
		l.Errorw("[GetOrderList] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryOrderListByPage error: %v", err.Error())
	}
	resp = &types.GetOrderListResponse{}
	resp.List = make([]types.Order, 0)
	tool.DeepCopy(&resp.List, list)
	resp.Total = total
	return
}
