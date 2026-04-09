package order

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	queue "github.com/perfect-panel/server/internal/jobs"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateOrderStatusInput struct {
	Body types.UpdateOrderStatusRequest
}

func UpdateOrderStatusHandler(deps Deps) func(context.Context, *UpdateOrderStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateOrderStatusInput) (*struct{}, error) {
		l := NewUpdateOrderStatusLogic(ctx, deps)
		if err := l.UpdateOrderStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateOrderStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update order status
func NewUpdateOrderStatusLogic(ctx context.Context, deps Deps) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusRequest) error {
	info, err := l.deps.OrderModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[UpdateOrderStatus] FindOne error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %v", err.Error())
	}

	if req.PaymentId != 0 {
		paymentMethod, err := l.deps.PaymentModel.FindOne(l.ctx, req.PaymentId)
		if err != nil {
			l.Error("[CreateOrder] PaymentMethod Not Found", logger.Field("error", err.Error()))
			return errors.Wrapf(xerr.NewErrCode(xerr.PaymentMethodNotFound), "PaymentMethod not found: %v", err.Error())
		}
		info.PaymentId = paymentMethod.Id
		info.Method = paymentMethod.Platform
	}
	if req.TradeNo != "" {
		info.TradeNo = req.TradeNo
	}

	err = l.deps.OrderModel.Transaction(l.ctx, func(db *gorm.DB) error {
		if err := l.deps.OrderModel.Update(l.ctx, info, db); err != nil {
			l.Errorw("[UpdateOrderStatus] Update error", logger.Field("error", err.Error()), logger.Field("OrderID", info.Id))
			return err
		}
		if err := l.deps.OrderModel.UpdateOrderStatus(l.ctx, info.OrderNo, req.Status, db); err != nil {
			return err
		}
		// If order status is 2, create user subscription
		if req.Status == 2 {
			payload := queue.ForthwithActivateOrderPayload{
				OrderNo: info.OrderNo,
			}
			p, _ := json.Marshal(payload)
			task := asynq.NewTask(queue.ForthwithActivateOrder, p)
			_, err = l.deps.Queue.EnqueueContext(l.ctx, task)
			if err != nil {
				l.Errorw("[UpdateOrderStatus] Enqueue error", logger.Field("error", err.Error()))
				return errors.Wrapf(xerr.NewErrCode(xerr.QueueEnqueueError), "Enqueue error: %v", err.Error())
			}
		}
		return nil
	})
	if err != nil {
		l.Errorw("[UpdateOrderStatus] Transaction error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Transaction error: %v", err.Error())
	}
	return nil
}
