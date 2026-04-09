package notify

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/config"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/platform/http/middleware"
	paymentModel "github.com/perfect-panel/server/models/payment"
	telegramsvc "github.com/perfect-panel/server/services/telegram"
	"github.com/stretchr/testify/require"
)

type stubPaymentModel struct {
	paymentModel.Model
	payment *paymentModel.Payment
	err     error
}

func (s stubPaymentModel) FindOneByPaymentToken(_ context.Context, _ string) (*paymentModel.Payment, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.payment, nil
}

func TestPhase5PaymentNotifyUnsupportedPlatformUsesPlainTextFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/payment/notify", withRequestContext(config.CtxKeyPlatform, "unsupported", PaymentNotifyHandler(Deps{})))

	req := httptest.NewRequest(http.MethodPost, "/payment/notify", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Equal(t, "unsupported platform", resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func TestPhase5NotifyMiddlewareLookupFailureKeepsEPayProtocolFailureShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	group := router.Group("/api/v1/notify")
	group.Use(middleware.NotifyMiddleware(&appruntime.Deps{
		PaymentModel: stubPaymentModel{err: errors.New("payment not found")},
	}))
	group.Any("/:platform/:token", PaymentNotifyHandler(Deps{}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notify/EPay/missing-token", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Equal(t, "invalid notification", resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func TestPhase5NotifyMiddlewareLookupFailureKeepsStripeProtocolFailureShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	group := router.Group("/api/v1/notify")
	group.Use(middleware.NotifyMiddleware(&appruntime.Deps{
		PaymentModel: stubPaymentModel{err: errors.New("payment not found")},
	}))
	group.Any("/:platform/:token", PaymentNotifyHandler(Deps{}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notify/Stripe/missing-token", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Empty(t, resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func TestPhase5StripeNotifyFailureKeepsProtocolAckShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST(
		"/payment/notify",
		withRequestContext(config.CtxKeyPlatform, "Stripe",
			withRequestContext(config.CtxKeyPayment, &paymentModel.Payment{Config: `{}`}, PaymentNotifyHandler(Deps{}))),
	)

	req := httptest.NewRequest(http.MethodPost, "/payment/notify", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Empty(t, resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func TestPhase5EPayInvalidSignatureKeepsProtocolFailureShape(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST(
		"/payment/notify",
		withRequestContext(config.CtxKeyPlatform, "EPay",
			withRequestContext(config.CtxKeyPayment, &paymentModel.Payment{
				Config: `{"pid":"pid","url":"https://example.com","key":"secret","type":"alipay"}`,
			}, PaymentNotifyHandler(Deps{}))),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/payment/notify?trade_no=123&out_trade_no=order-1&trade_status=TRADE_SUCCESS&sign=fake",
		nil,
	)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Equal(t, "invalid notification", resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func TestPhase5TelegramSecretMismatchKeepsEmptyAck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	telegramCfg := config.Config{
		Telegram: config.Telegram{BotToken: "bot-token"},
	}
	telegramsvc.RegisterTelegramHandlers(router, telegramsvc.Deps{
		Config: &telegramCfg,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/telegram/webhook?secret=wrong", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Empty(t, resp.Body.String())
	require.NotContains(t, resp.Header().Get("Content-Type"), "application/problem+json")
}

func withRequestContext(key config.CtxKey, value interface{}, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), key, value)
		c.Request = c.Request.WithContext(ctx)
		next(c)
	}
}
