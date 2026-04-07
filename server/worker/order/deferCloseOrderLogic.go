package orderLogic

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/services/user/order"
	internal "github.com/perfect-panel/server/types"
	"github.com/perfect-panel/server/worker/spec"
)

type DeferCloseOrderLogic struct {
	deps Deps
}

func NewDeferCloseOrderLogic(deps Deps) *DeferCloseOrderLogic {
	return &DeferCloseOrderLogic{
		deps: deps,
	}
}

func (l *DeferCloseOrderLogic) ProcessTask(ctx context.Context, task *asynq.Task) error {
	payload := spec.DeferCloseOrderPayload{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.WithContext(ctx).Error("[DeferCloseOrderLogic] Unmarshal payload failed",
			logger.Field("error", err.Error()),
			logger.Field("payload", string(task.Payload())),
		)
		return nil
	}

	orderDeps := order.Deps{
		OrderModel:     l.deps.OrderModel,
		PaymentModel:   l.deps.PaymentModel,
		SubscribeModel: l.deps.SubscribeModel,
		UserModel:      l.deps.UserModel,
		CouponModel:    l.deps.CouponModel,
		DB:             l.deps.DB,
		Queue:          l.deps.Queue,
		Config:         l.deps.Config,
	}
	err := order.NewCloseOrderLogic(ctx, orderDeps).CloseOrder(&internal.CloseOrderRequest{
		OrderNo: payload.OrderNo,
	})
	count, ok := asynq.GetRetryCount(ctx)
	if !ok {
		return nil
	}
	if err != nil && count < 3 {
		return err
	}
	return nil
}
