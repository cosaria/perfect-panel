package order

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type QueryOrderDetailInput struct {
	types.QueryOrderDetailRequest
}

type QueryOrderDetailOutput struct {
	Body *types.OrderDetail
}

func QueryOrderDetailHandler(deps Deps) func(context.Context, *QueryOrderDetailInput) (*QueryOrderDetailOutput, error) {
	return func(ctx context.Context, input *QueryOrderDetailInput) (*QueryOrderDetailOutput, error) {
		l := NewQueryOrderDetailLogic(ctx, deps)
		resp, err := l.QueryOrderDetail(&input.QueryOrderDetailRequest)
		if err != nil {
			return nil, err
		}
		return &QueryOrderDetailOutput{Body: resp}, nil
	}
}

type QueryOrderDetailLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get order
func NewQueryOrderDetailLogic(ctx context.Context, deps Deps) *QueryOrderDetailLogic {
	return &QueryOrderDetailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryOrderDetailLogic) QueryOrderDetail(req *types.QueryOrderDetailRequest) (resp *types.OrderDetail, err error) {
	orderInfo, err := l.deps.OrderModel.FindOneDetailsByOrderNo(l.ctx, req.OrderNo)
	if err != nil {
		l.Errorw("[QueryOrderDetail] Database query error", logger.Field("error", err.Error()), logger.Field("order_no", req.OrderNo))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find order error: %v", err.Error())
	}
	resp = &types.OrderDetail{}
	tool.DeepCopy(resp, orderInfo)
	// Prevent commission amount leakage
	resp.Commission = 0
	return
}
