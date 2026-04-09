package middleware

import (
	"context"
	"net/http"

	serverconfig "github.com/perfect-panel/server/config"

	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/platform/payment"
)

type PaymentParams struct {
	Platform string `uri:"platform"`
	Token    string `uri:"token"`
}

func NotifyMiddleware(runtimeDeps *appruntime.Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var params PaymentParams
		// Get platform and token from uri
		if err := c.ShouldBindUri(&params); err != nil {
			writeNotifyProtocolFailure(c, params.Platform)
			c.Abort()
			return
		}
		paymentConfig, err := runtimeDeps.PaymentModel.FindOneByPaymentToken(ctx, params.Token)
		if err != nil {
			writeNotifyProtocolFailure(c, params.Platform)
			c.Abort()
			return
		}
		ctx = context.WithValue(ctx, serverconfig.CtxKeyPlatform, paymentConfig.Platform)
		ctx = context.WithValue(ctx, serverconfig.CtxKeyPayment, paymentConfig)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func writeNotifyProtocolFailure(c *gin.Context, platform string) {
	switch payment.ParsePlatform(platform) {
	case payment.Stripe:
		c.Status(http.StatusBadRequest)
	case payment.EPay, payment.CryptoSaaS, payment.AlipayF2F:
		c.String(http.StatusBadRequest, "%s", "invalid notification")
	default:
		c.String(http.StatusBadRequest, "%s", "unsupported platform")
	}
}
