package notify

import (
	"net/http"

	"github.com/perfect-panel/server/config"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/payment"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

// PaymentNotifyHandler Payment Notify
func PaymentNotifyHandler(deps Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		platform, ok := c.Request.Context().Value(config.CtxKeyPlatform).(string)
		if !ok {
			logger.WithContext(c.Request.Context()).Errorf("platform not found")
			writePlainText(c, http.StatusBadRequest, "unsupported platform")
			return
		}

		switch payment.ParsePlatform(platform) {
		case payment.EPay, payment.CryptoSaaS:
			req := &types.EPayNotifyRequest{}
			if err := c.ShouldBind(req); err != nil {
				writePlainText(c, http.StatusBadRequest, "invalid notification")
				return
			}
			l := NewEPayNotifyLogic(c, deps)
			if err := l.EPayNotify(req); err != nil {
				logger.WithContext(c.Request.Context()).Errorf("EPayNotify failed: %v", err.Error())
				writeProtocolFailure(c, err, http.StatusBadRequest, "invalid notification", "failed", false)
				return
			}
			writePlainText(c, http.StatusOK, "success")
		case payment.Stripe:
			l := NewStripeNotifyLogic(c.Request.Context(), deps)
			if err := l.StripeNotify(c.Request, c.Writer); err != nil {
				writeProtocolFailure(c, err, http.StatusBadRequest, "", "", true)
				return
			}
			writeEmptyStatus(c, http.StatusOK)

		case payment.AlipayF2F:
			l := NewAlipayNotifyLogic(c.Request.Context(), deps)
			if err := l.AlipayNotify(c.Request); err != nil {
				writeProtocolFailure(c, err, http.StatusBadRequest, "invalid notification", "failed", false)
				return
			}
			writePlainText(c, http.StatusOK, "success")

		default:
			logger.WithContext(c.Request.Context()).Errorf("platform %s not support", platform)
			writePlainText(c, http.StatusBadRequest, "unsupported platform")
		}
	}
}
