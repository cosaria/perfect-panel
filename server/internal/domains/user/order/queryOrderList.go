package order

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type QueryOrderListInput struct {
	types.QueryOrderListRequest
}

type QueryOrderListOutput struct {
	Body *types.QueryOrderListResponse
}

func QueryOrderListHandler(deps Deps) func(context.Context, *QueryOrderListInput) (*QueryOrderListOutput, error) {
	return func(ctx context.Context, input *QueryOrderListInput) (*QueryOrderListOutput, error) {
		l := NewQueryOrderListLogic(ctx, deps)
		resp, err := l.QueryOrderList(&input.QueryOrderListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryOrderListOutput{Body: resp}, nil
	}
}

type QueryOrderListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get order list
func NewQueryOrderListLogic(ctx context.Context, deps Deps) *QueryOrderListLogic {
	return &QueryOrderListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryOrderListLogic) QueryOrderList(req *types.QueryOrderListRequest) (resp *types.QueryOrderListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	total, data, err := l.deps.OrderModel.QueryOrderListByPage(l.ctx, req.Page, req.Size, 0, u.Id, 0, "")
	if err != nil {
		l.Errorw("[QueryOrderListLogic] Query order list failed", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query order list failed")
	}
	resp = &types.QueryOrderListResponse{
		Total: total,
		List:  make([]types.OrderDetail, 0),
	}
	for _, item := range data {
		var orderInfo types.OrderDetail
		tool.DeepCopy(&orderInfo, item)
		// Prevent commission amount leakage
		orderInfo.Commission = 0
		resp.List = append(resp.List, orderInfo)
	}

	return
}
