package order

import (
	"context"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateOrderInput struct {
	Body types.CreateOrderRequest
}

func CreateOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateOrderInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateOrderInput) (*struct{}, error) {
		l := NewCreateOrderLogic(ctx, svcCtx)
		if err := l.CreateOrder(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateOrderLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create order
func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) error {
	paymentMethod, err := l.svcCtx.PaymentModel.FindOne(l.ctx, req.PaymentId)
	if err != nil {
		l.Logger.Error("[CreateOrder] PaymentMethod Not Found", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.PaymentMethodNotFound), "PaymentMethod not found: %v", err.Error())
	}

	err = l.svcCtx.OrderModel.Insert(l.ctx, &order.Order{
		UserId:         req.UserId,
		OrderNo:        tool.GenerateTradeNo(),
		Type:           req.Type,
		Quantity:       req.Quantity,
		Price:          req.Price,
		Amount:         req.Amount,
		Discount:       req.Discount,
		Coupon:         req.Coupon,
		CouponDiscount: req.CouponDiscount,
		PaymentId:      req.PaymentId,
		Method:         paymentMethod.Token,
		FeeAmount:      req.FeeAmount,
		TradeNo:        req.TradeNo,
		Status:         req.Status,
		SubscribeId:    req.SubscribeId,
	})
	if err != nil {
		l.Logger.Error("[CreateOrder] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Insert error: %v", err.Error())
	}
	return nil
}
