package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/jobs/spec"
	"github.com/perfect-panel/server/internal/platform/payment/alipay"
	"github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type AlipayNotifyLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Alipay notify
func NewAlipayNotifyLogic(ctx context.Context, deps Deps) *AlipayNotifyLogic {
	return &AlipayNotifyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *AlipayNotifyLogic) AlipayNotify(r *http.Request) error {
	data, ok := l.ctx.Value(config.CtxKeyPayment).(*payment.Payment)
	if !ok {
		return fmt.Errorf("payment config not found")
	}
	var config payment.AlipayF2FConfig
	if err := json.Unmarshal([]byte(data.Config), &config); err != nil {
		l.Error("[AlipayNotify] Unmarshal config failed", logger.Field("error", err.Error()))
		return err
	}
	client := alipay.NewClient(alipay.Config{
		AppId:       config.AppId,
		PrivateKey:  config.PrivateKey,
		PublicKey:   config.PublicKey,
		InvoiceName: config.InvoiceName,
		NotifyURL:   data.Domain + "/api/v1/payment/alipay/notify",
	})
	notify, err := client.DecodeNotification(r.Form)
	if err != nil {
		l.Error("[AlipayNotify] Decode notification failed", logger.Field("error", err.Error()))
		return markInvalidNotification(err)
	}
	if notify.Status == alipay.Success {
		orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, notify.OrderNo)
		if err != nil {
			l.Error("[AlipayNotify] Find order failed", logger.Field("error", err.Error()), logger.Field("orderNo", notify.OrderNo))
			return errors.Wrapf(xerr.NewErrCode(xerr.OrderNotExist), "order not exist: %v", notify.OrderNo)
		}

		if orderInfo.Status == 5 {
			return nil
		}

		// Update order status
		err = l.deps.OrderModel.UpdateOrderStatus(l.ctx, notify.OrderNo, 2)
		if err != nil {
			l.Error("[AlipayNotify] Update order status failed", logger.Field("error", err.Error()), logger.Field("orderNo", notify.OrderNo))
			return err
		}
		l.Info("[AlipayNotify] Notify status success", logger.Field("orderNo", notify.OrderNo))
		payload := spec.ForthwithActivateOrderPayload{
			OrderNo: notify.OrderNo,
		}
		bytes, err := json.Marshal(&payload)
		if err != nil {
			l.Error("[AlipayNotify] Marshal payload failed", logger.Field("error", err.Error()))
			return err
		}
		task := asynq.NewTask(spec.ForthwithActivateOrder, bytes)
		taskInfo, err := l.deps.Queue.EnqueueContext(l.ctx, task)
		if err != nil {
			l.Error("[AlipayNotify] Enqueue task failed", logger.Field("error", err.Error()))
			return err
		}
		l.Info("[AlipayNotify] Enqueue task success", logger.Field("taskInfo", taskInfo))
	} else {
		l.Error("[AlipayNotify] Notify status failed", logger.Field("status", string(notify.Status)))
	}
	return nil
}
