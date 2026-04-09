package order

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	queue "github.com/perfect-panel/server/internal/jobs/spec"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type RechargeInput struct {
	Body types.RechargeOrderRequest
}

type RechargeOutput struct {
	Body *types.RechargeOrderResponse
}

func RechargeHandler(deps Deps) func(context.Context, *RechargeInput) (*RechargeOutput, error) {
	return func(ctx context.Context, input *RechargeInput) (*RechargeOutput, error) {
		l := NewRechargeLogic(ctx, deps)
		resp, err := l.Recharge(&input.Body)
		if err != nil {
			return nil, err
		}
		return &RechargeOutput{Body: resp}, nil
	}
}

type RechargeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Recharge
func NewRechargeLogic(ctx context.Context, deps Deps) *RechargeLogic {
	return &RechargeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *RechargeLogic) Recharge(req *types.RechargeOrderRequest) (resp *types.RechargeOrderResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	// Validate recharge amount
	if req.Amount <= 0 {
		l.Errorw("[Recharge] Invalid recharge amount", logger.Field("amount", req.Amount), logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "recharge amount must be greater than 0")
	}

	if req.Amount > MaxRechargeAmount {
		l.Errorw("[Recharge] Recharge amount exceeds maximum limit",
			logger.Field("amount", req.Amount),
			logger.Field("max", MaxRechargeAmount),
			logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "recharge amount exceeds maximum limit")
	}

	// find payment method
	payment, err := l.deps.PaymentModel.FindOne(l.ctx, req.Payment)
	if err != nil {
		l.Errorw("[Recharge] Database query error", logger.Field("error", err.Error()), logger.Field("payment", req.Payment))
		return nil, errors.Wrapf(err, "find payment error: %v", err.Error())
	}
	// Calculate the handling fee
	feeAmount := calculateFee(req.Amount, payment)
	totalAmount := req.Amount + feeAmount

	// Validate total amount after adding fee
	if totalAmount > MaxOrderAmount {
		l.Errorw("[Recharge] Total amount exceeds maximum limit after fee",
			logger.Field("amount", totalAmount),
			logger.Field("max", MaxOrderAmount),
			logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "total amount exceeds maximum limit")
	}

	// query user is new purchase or renewal
	isNew, err := l.deps.OrderModel.IsUserEligibleForNewOrder(l.ctx, u.Id)
	if err != nil {
		l.Errorw("[Recharge] Database query error", logger.Field("error", err.Error()), logger.Field("user_id", u.Id))
		return nil, errors.Wrapf(err, "query user error: %v", err.Error())
	}
	orderInfo := order.Order{
		UserId:    u.Id,
		OrderNo:   tool.GenerateTradeNo(),
		Type:      4,
		Price:     req.Amount,
		Amount:    totalAmount,
		FeeAmount: feeAmount,
		PaymentId: payment.Id,
		Method:    payment.Platform,
		Status:    1,
		IsNew:     isNew,
	}
	err = l.deps.OrderModel.Insert(l.ctx, &orderInfo)
	if err != nil {
		l.Errorw("[Recharge] Database insert error", logger.Field("error", err.Error()), logger.Field("order", orderInfo))
		return nil, errors.Wrapf(err, "insert order error: %v", err.Error())
	}
	// Deferred task
	payload := queue.DeferCloseOrderPayload{
		OrderNo: orderInfo.OrderNo,
	}
	val, err := json.Marshal(payload)
	if err != nil {
		l.Errorw("[Recharge] Marshal payload error", logger.Field("error", err.Error()), logger.Field("payload", payload))
	}
	task := asynq.NewTask(queue.DeferCloseOrder, val, asynq.MaxRetry(3))
	taskInfo, err := l.deps.Queue.Enqueue(task, asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
	if err != nil {
		l.Errorw("[Recharge] Enqueue task error", logger.Field("error", err.Error()), logger.Field("task", task))
	} else {
		l.Infow("[Recharge] Enqueue task success", logger.Field("TaskID", taskInfo.ID))
	}
	return &types.RechargeOrderResponse{
		OrderNo: orderInfo.OrderNo,
	}, nil
}
