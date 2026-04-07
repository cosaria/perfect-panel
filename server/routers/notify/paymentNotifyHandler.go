package notify

import (
	"fmt"
	"net/http"

	"github.com/perfect-panel/server/config"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/payment"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/services/notify"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

// PaymentNotifyHandler Payment Notify
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		platform, ok := c.Request.Context().Value(config.CtxKeyPlatform).(string)
		if !ok {
			logger.WithContext(c.Request.Context()).Errorf("platform not found")
			response.HttpResult(c, nil, fmt.Errorf("platform not found"))
			return
		}

		switch payment.ParsePlatform(platform) {
		case payment.EPay, payment.CryptoSaaS:
			req := &types.EPayNotifyRequest{}
			if err := c.ShouldBind(req); err != nil {
				response.HttpResult(c, nil, err)
				return
			}
			l := notify.NewEPayNotifyLogic(c, svcCtx)
			if err := l.EPayNotify(req); err != nil {
				logger.WithContext(c.Request.Context()).Errorf("EPayNotify failed: %v", err.Error())
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			c.String(http.StatusOK, "%s", "success")
		case payment.Stripe:
			l := notify.NewStripeNotifyLogic(c.Request.Context(), svcCtx)
			if err := l.StripeNotify(c.Request, c.Writer); err != nil {
				response.HttpResult(c, nil, err)
				return
			}
			response.HttpResult(c, nil, nil)

		case payment.AlipayF2F:
			l := notify.NewAlipayNotifyLogic(c.Request.Context(), svcCtx)
			if err := l.AlipayNotify(c.Request); err != nil {
				response.HttpResult(c, nil, err)
				return
			}
			// Return success to alipay
			c.String(http.StatusOK, "%s", "success")

		default:
			logger.WithContext(c.Request.Context()).Errorf("platform %s not support", platform)
		}
	}
}
