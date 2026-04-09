package notify

import (
	"encoding/json"
	stderrors "errors"
	"net/url"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/payment/epay"
	"github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/support/logger"

	queueType "github.com/perfect-panel/server/internal/jobs"
)

type EPayNotifyLogic struct {
	logger.Logger
	ctx  *gin.Context
	deps Deps
}

// EPay notify
func NewEPayNotifyLogic(ctx *gin.Context, deps Deps) *EPayNotifyLogic {
	return &EPayNotifyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *EPayNotifyLogic) EPayNotify(req *types.EPayNotifyRequest) error {
	rawPayload := l.ctx.Request.URL.RawQuery
	idempotencyKey := callbackHash("epay", rawPayload)

	// Find payment config
	data, ok := l.ctx.Request.Context().Value(config.CtxKeyPayment).(*payment.Payment)
	if !ok {
		l.Error("[EPayNotify] Payment not found in context")
		recordExternalTrust(l.ctx, l.deps, &system.ExternalTrustEvent{
			EntryPoint:      "payment_notify",
			IdempotencyKey:  idempotencyKey,
			AuthStatus:      "failed",
			ProcessingState: "rejected",
			FailureReason:   "payment config not found",
			RawPayload:      rawPayload,
		})
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "payment config not found")
	}
	l.Infof("[EPayNotify] Payment config: %+v", data)
	recordTrust := func(authStatus string, state string, failure string) {
		recordExternalTrust(l.ctx, l.deps, &system.ExternalTrustEvent{
			EntryPoint:      "payment_notify",
			Credential:      data.Token,
			IdempotencyKey:  idempotencyKey,
			AuthStatus:      authStatus,
			ProcessingState: state,
			FailureReason:   failure,
			RawPayload:      rawPayload,
		})
	}

	var config payment.EPayConfig
	if err := json.Unmarshal([]byte(data.Config), &config); err != nil {
		recordTrust("failed", "rejected", err.Error())
		l.Errorw("[EPayNotify] Unmarshal config failed", logger.Field("error", err.Error()))
		return err
	}
	// Verify sign
	client := epay.NewClient(config.Pid, config.Url, config.Key, config.Type)
	if !client.VerifySign(urlParamsToMap(l.ctx.Request.URL.RawQuery)) && !l.deps.debugEnabled() {
		recordTrust("failed", "rejected", "verify sign failed")
		l.Error("[EPayNotify] Verify sign failed")
		return markInvalidNotification(stderrors.New("verify sign failed"))
	}
	if req.TradeStatus != "TRADE_SUCCESS" {
		recordTrust("verified", "ignored", "")
		l.Error("[EPayNotify] Trade status is not success", logger.Field("orderNo", req.OutTradeNo), logger.Field("tradeStatus", req.TradeStatus))
		return nil
	}
	orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, req.OutTradeNo)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		l.Error("[EPayNotify] Find order failed", logger.Field("error", err.Error()), logger.Field("orderNo", req.OutTradeNo))
		return errors.Wrapf(xerr.NewErrCode(xerr.OrderNotExist), "order not exist: %v", req.OutTradeNo)
	}
	decision, err := recordPaymentCallbackAttempt(l.ctx, l.deps, orderInfo.PaymentId, "epay", idempotencyKey, rawPayload)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		return err
	}
	if !decision.Accepted {
		recordTrust("verified", "duplicate", "")
		return nil
	}
	if orderInfo.Status == 5 {
		markPaymentCallbackProcessed(l.ctx, l.deps, decision.CallbackID, "processed")
		recordTrust("verified", "processed", "")
		return nil
	}
	// Update order status
	err = l.deps.OrderModel.UpdateOrderStatus(l.ctx, req.OutTradeNo, 2)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		l.Error("[EPayNotify] Update order status failed", logger.Field("error", err.Error()), logger.Field("orderNo", req.OutTradeNo))
		return err
	}
	// Create activate order task
	payload := queueType.ForthwithActivateOrderPayload{
		OrderNo: req.OutTradeNo,
	}
	bytes, err := json.Marshal(&payload)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		l.Error("[EPayNotify] Marshal payload failed", logger.Field("error", err.Error()))
		return err
	}
	task := asynq.NewTask(queueType.ForthwithActivateOrder, bytes)
	taskInfo, err := l.deps.Queue.EnqueueContext(l.ctx, task)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		l.Error("[EPayNotify] Enqueue task failed", logger.Field("error", err.Error()))
		return err
	}
	markPaymentCallbackProcessed(l.ctx, l.deps, decision.CallbackID, "processed")
	recordTrust("verified", "processed", "")
	l.Info("[EPayNotify] Enqueue task success", logger.Field("taskInfo", taskInfo))
	return nil
}

func urlParamsToMap(query string) map[string]string {
	params := make(map[string]string)
	values, _ := url.ParseQuery(query)
	for k, v := range values {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}
