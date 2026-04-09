package notify

import (
	"encoding/json"
	stderrors "errors"
	"net/url"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/payment"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/payment/epay"

	queueType "github.com/perfect-panel/server/worker"
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

	// Find payment config
	data, ok := l.ctx.Request.Context().Value(config.CtxKeyPayment).(*payment.Payment)
	if !ok {
		l.Error("[EPayNotify] Payment not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "payment config not found")
	}
	l.Infof("[EPayNotify] Payment config: %+v", data)

	var config payment.EPayConfig
	if err := json.Unmarshal([]byte(data.Config), &config); err != nil {
		l.Errorw("[EPayNotify] Unmarshal config failed", logger.Field("error", err.Error()))
		return err
	}
	// Verify sign
	client := epay.NewClient(config.Pid, config.Url, config.Key, config.Type)
	if !client.VerifySign(urlParamsToMap(l.ctx.Request.URL.RawQuery)) && !l.deps.debugEnabled() {
		l.Error("[EPayNotify] Verify sign failed")
		return markInvalidNotification(stderrors.New("verify sign failed"))
	}
	if req.TradeStatus != "TRADE_SUCCESS" {
		l.Error("[EPayNotify] Trade status is not success", logger.Field("orderNo", req.OutTradeNo), logger.Field("tradeStatus", req.TradeStatus))
		return nil
	}
	orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, req.OutTradeNo)
	if err != nil {
		l.Error("[EPayNotify] Find order failed", logger.Field("error", err.Error()), logger.Field("orderNo", req.OutTradeNo))
		return errors.Wrapf(xerr.NewErrCode(xerr.OrderNotExist), "order not exist: %v", req.OutTradeNo)
	}
	if orderInfo.Status == 5 {
		return nil
	}
	// Update order status
	err = l.deps.OrderModel.UpdateOrderStatus(l.ctx, req.OutTradeNo, 2)
	if err != nil {
		l.Error("[EPayNotify] Update order status failed", logger.Field("error", err.Error()), logger.Field("orderNo", req.OutTradeNo))
		return err
	}
	// Create activate order task
	payload := queueType.ForthwithActivateOrderPayload{
		OrderNo: req.OutTradeNo,
	}
	bytes, err := json.Marshal(&payload)
	if err != nil {
		l.Error("[EPayNotify] Marshal payload failed", logger.Field("error", err.Error()))
		return err
	}
	task := asynq.NewTask(queueType.ForthwithActivateOrder, bytes)
	taskInfo, err := l.deps.Queue.EnqueueContext(l.ctx, task)
	if err != nil {
		l.Error("[EPayNotify] Enqueue task failed", logger.Field("error", err.Error()))
		return err
	}
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
