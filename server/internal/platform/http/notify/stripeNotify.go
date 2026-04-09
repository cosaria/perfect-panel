package notify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/jobs/spec"
	"github.com/perfect-panel/server/internal/platform/payment/stripe"
	"github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type StripeNotifyLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewStripeNotifyLogic Stripe notify
func NewStripeNotifyLogic(ctx context.Context, deps Deps) *StripeNotifyLogic {
	return &StripeNotifyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *StripeNotifyLogic) StripeNotify(r *http.Request, w http.ResponseWriter) error {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	idempotencyKey := callbackHash("stripe", string(payload))
	if err != nil {
		recordExternalTrust(l.ctx, l.deps, &system.ExternalTrustEvent{
			EntryPoint:      "payment_notify",
			IdempotencyKey:  idempotencyKey,
			AuthStatus:      "failed",
			ProcessingState: "rejected",
			FailureReason:   err.Error(),
		})
		l.Errorw("[StripeNotify] error", logger.Field("errors", err.Error()))
		return markInvalidNotification(err)
	}
	signature := r.Header.Get("Stripe-Signature")
	stripeConfig, ok := l.ctx.Value(config.CtxKeyPayment).(*payment.Payment)
	if !ok {
		recordExternalTrust(l.ctx, l.deps, &system.ExternalTrustEvent{
			EntryPoint:      "payment_notify",
			IdempotencyKey:  idempotencyKey,
			AuthStatus:      "failed",
			ProcessingState: "rejected",
			FailureReason:   "payment config not found",
			RawPayload:      string(payload),
		})
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "payment config not found")
	}
	recordTrust := func(authStatus string, state string, failure string) {
		recordExternalTrust(l.ctx, l.deps, &system.ExternalTrustEvent{
			EntryPoint:      "payment_notify",
			Credential:      stripeConfig.Token,
			IdempotencyKey:  idempotencyKey,
			AuthStatus:      authStatus,
			ProcessingState: state,
			FailureReason:   failure,
			RawPayload:      string(payload),
		})
	}
	config := payment.StripeConfig{}
	if err := json.Unmarshal([]byte(stripeConfig.Config), &config); err != nil {
		recordTrust("failed", "rejected", err.Error())
		return err
	}
	client := stripe.NewClient(stripe.Config{
		PublicKey:     config.PublicKey,
		SecretKey:     config.SecretKey,
		WebhookSecret: config.WebhookSecret,
	})

	notify, err := client.ParseNotify(payload, signature)
	if err != nil {
		recordTrust("failed", "rejected", err.Error())
		l.Errorw("[StripeNotify] error", logger.Field("errors", err.Error()))
		return markInvalidNotification(err)
	}
	orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, notify.OrderNo)
	if err != nil {
		recordTrust("verified", "failed", err.Error())
		l.Error("[StripeNotify] Find order failed", logger.Field("error", err.Error()), logger.Field("orderNo", notify.OrderNo))
		return errors.Wrapf(xerr.NewErrCode(xerr.OrderNotExist), "order not exist: %v", notify.OrderNo)
	}
	if notify.EventType == "payment_intent.succeeded" {
		decision, err := recordPaymentCallbackAttempt(l.ctx, l.deps, orderInfo.PaymentId, "stripe", idempotencyKey, string(payload))
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
		// update order status
		err = l.deps.OrderModel.UpdateOrderStatus(l.ctx, notify.OrderNo, 2)
		if err != nil {
			recordTrust("verified", "failed", err.Error())
			return err
		}
		// create ActivateOrder task
		payload := spec.ForthwithActivateOrderPayload{
			OrderNo: notify.OrderNo,
		}
		bytes, err := json.Marshal(payload)
		if err != nil {
			recordTrust("verified", "failed", err.Error())
			l.Errorw("[StripeNotify] Marshal error", logger.Field("errors", err.Error()), logger.Field("payload", payload))
			return err
		}
		task := asynq.NewTask(spec.ForthwithActivateOrder, bytes)
		_, err = l.deps.Queue.Enqueue(task)
		if err != nil {
			recordTrust("verified", "failed", err.Error())
			l.Errorw("[StripeNotify] Enqueue error", logger.Field("errors", err.Error()))
			return err
		}
		markPaymentCallbackProcessed(l.ctx, l.deps, decision.CallbackID, "processed")
		recordTrust("verified", "processed", "")
		l.Infow("[StripeNotify] success", logger.Field("orderNo", notify.OrderNo))
	}
	return nil
}
