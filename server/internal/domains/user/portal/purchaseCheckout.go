package portal

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/domains/common/report"
	queueType "github.com/perfect-panel/server/internal/jobs"
	"github.com/perfect-panel/server/internal/platform/http/types"
	paymentPlatform "github.com/perfect-panel/server/internal/platform/payment"
	"github.com/perfect-panel/server/internal/platform/payment/alipay"
	"github.com/perfect-panel/server/internal/platform/payment/epay"
	"github.com/perfect-panel/server/internal/platform/payment/exchangeRate"
	"github.com/perfect-panel/server/internal/platform/payment/stripe"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/persistence/order"
	"github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PurchaseCheckoutInput struct {
	Body types.CheckoutOrderRequest
}

type PurchaseCheckoutOutput struct {
	Body *types.CheckoutOrderResponse
}

func PurchaseCheckoutHandler(deps Deps) func(context.Context, *PurchaseCheckoutInput) (*PurchaseCheckoutOutput, error) {
	return func(ctx context.Context, input *PurchaseCheckoutInput) (*PurchaseCheckoutOutput, error) {
		l := NewPurchaseCheckoutLogic(ctx, deps)
		resp, err := l.PurchaseCheckout(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PurchaseCheckoutOutput{Body: resp}, nil
	}
}

type PurchaseCheckoutLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewPurchaseCheckoutLogic creates a new instance of PurchaseCheckoutLogic
// for handling purchase checkout operations across different payment platforms
func NewPurchaseCheckoutLogic(ctx context.Context, deps Deps) *PurchaseCheckoutLogic {
	return &PurchaseCheckoutLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

// PurchaseCheckout processes the checkout for an order using the specified payment method
// It validates the order, retrieves payment configuration, and routes to the appropriate payment handler
func (l *PurchaseCheckoutLogic) PurchaseCheckout(req *types.CheckoutOrderRequest) (resp *types.CheckoutOrderResponse, err error) {

	// Validate and retrieve order information
	orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, req.OrderNo)
	if err != nil {
		l.Error("[PurchaseCheckout] Find order failed", logger.Field("error", err.Error()), logger.Field("orderNo", req.OrderNo))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.OrderNotExist), "order not exist: %v", req.OrderNo)
	}

	// Verify order is in pending payment status (status = 1)
	if orderInfo.Status != 1 {
		l.Error("[PurchaseCheckout] Order status error", logger.Field("status", orderInfo.Status))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.OrderStatusError), "order status error: %v", orderInfo.Status)
	}

	// Retrieve payment method configuration
	paymentConfig, err := l.deps.PaymentModel.FindOne(l.ctx, orderInfo.PaymentId)
	if err != nil {
		l.Error("[PurchaseCheckout] Database query error", logger.Field("error", err.Error()), logger.Field("payment", orderInfo.Method))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find payment method error: %v", err.Error())
	}
	// Route to appropriate payment handler based on payment platform
	switch paymentPlatform.ParsePlatform(orderInfo.Method) {
	case paymentPlatform.EPay:
		// Process EPay payment - generates payment URL for redirect
		url, err := l.epayPayment(paymentConfig, orderInfo, req.ReturnUrl)
		if err != nil {
			l.Error("[PurchaseCheckout] epay error", logger.Field("error", err.Error()))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "epayPayment error: %v", err.Error())
		}
		resp = &types.CheckoutOrderResponse{
			CheckoutUrl: url,
			Type:        "url", // Client should redirect to URL
		}

	case paymentPlatform.Stripe:
		// Process Stripe payment - creates payment sheet for client-side processing
		stripePayment, err := l.stripePayment(paymentConfig.Config, orderInfo, "")
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "stripePayment error: %v", err.Error())
		}
		resp = &types.CheckoutOrderResponse{
			Type:   "stripe", // Client should use Stripe SDK
			Stripe: stripePayment,
		}

	case paymentPlatform.AlipayF2F:
		// Process Alipay Face-to-Face payment - generates QR code
		url, err := l.alipayF2fPayment(paymentConfig, orderInfo)
		if err != nil {
			l.Errorw("[PurchaseCheckout] alipayF2fPayment error", logger.Field("error", err.Error()))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "alipayF2fPayment error: %v", err.Error())
		}
		resp = &types.CheckoutOrderResponse{
			Type:        "qr", // Client should display QR code
			CheckoutUrl: url,
		}

	case paymentPlatform.CryptoSaaS:
		// Process EPay payment - generates payment URL for redirect
		url, err := l.CryptoSaaSPayment(paymentConfig, orderInfo, req.ReturnUrl)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "epayPayment error: %v", err.Error())
		}
		resp = &types.CheckoutOrderResponse{
			CheckoutUrl: url,
			Type:        "url", // Client should redirect to URL
		}

	case paymentPlatform.Balance:
		// Process balance payment - validate user and process payment immediately
		if orderInfo.UserId == 0 {
			l.Errorw("[PurchaseCheckout] user not found", logger.Field("userId", orderInfo.UserId))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserNotExist), "user not found")
		}

		// Retrieve user information for balance validation
		userInfo, err := l.deps.UserModel.FindOne(l.ctx, orderInfo.UserId)
		if err != nil {
			l.Errorw("[PurchaseCheckout] FindOne User error", logger.Field("error", err.Error()))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %s", err.Error())
		}

		// Process balance payment with gift amount priority logic
		if err = l.balancePayment(userInfo, orderInfo); err != nil {
			return nil, err
		}

		resp = &types.CheckoutOrderResponse{
			Type: "balance", // Payment completed immediately
		}

	default:
		l.Errorw("[PurchaseCheckout] payment method not found", logger.Field("method", orderInfo.Method))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "payment method not found")
	}
	return
}

// alipayF2fPayment processes Alipay Face-to-Face payment by generating a QR code
// It handles currency conversion and creates a pre-payment trade for QR code scanning
func (l *PurchaseCheckoutLogic) alipayF2fPayment(pay *payment.Payment, info *order.Order) (string, error) {
	// Parse Alipay F2F configuration from payment settings
	f2FConfig := &payment.AlipayF2FConfig{}
	if err := f2FConfig.Unmarshal([]byte(pay.Config)); err != nil {
		l.Errorw("[PurchaseCheckout] Unmarshal Alipay config error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal error: %s", err.Error())
	}

	// Build notification URL for payment status callbacks
	notifyUrl := ""
	if pay.Domain != "" {
		notifyUrl = pay.Domain + "/api/v1/notify/" + pay.Platform + "/" + pay.Token
	} else {
		host, ok := l.ctx.Value(config.CtxKeyRequestHost).(string)
		if !ok {
			host = l.deps.Config.Host
		}
		notifyUrl = "https://" + host + "/api/v1/notify/" + pay.Platform + "/" + pay.Token
	}

	// Initialize Alipay client with configuration
	client := alipay.NewClient(alipay.Config{
		AppId:       f2FConfig.AppId,
		PrivateKey:  f2FConfig.PrivateKey,
		PublicKey:   f2FConfig.PublicKey,
		InvoiceName: f2FConfig.InvoiceName,
		NotifyURL:   notifyUrl,
	})

	// Convert order amount to CNY using current exchange rate
	amount, err := l.queryExchangeRate("CNY", info.Amount)
	if err != nil {
		l.Errorw("[PurchaseCheckout] queryExchangeRate error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "queryExchangeRate error: %s", err.Error())
	}
	convertAmount := int64(amount * 100) // Convert to cents for API

	// Create pre-payment trade and generate QR code
	QRCode, err := client.PreCreateTrade(l.ctx, alipay.Order{
		OrderNo: info.OrderNo,
		Amount:  convertAmount,
	})
	if err != nil {
		l.Errorw("[PurchaseCheckout] PreCreateTrade error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "PreCreateTrade error: %s", err.Error())
	}
	return QRCode, nil
}

// stripePayment processes Stripe payment by creating a payment sheet
// It supports various payment methods including WeChat Pay and Alipay through Stripe
func (l *PurchaseCheckoutLogic) stripePayment(config string, info *order.Order, identifier string) (*types.StripePayment, error) {
	// Parse Stripe configuration from payment settings
	stripeConfig := &payment.StripeConfig{}

	if err := stripeConfig.Unmarshal([]byte(config)); err != nil {
		l.Errorw("[PurchaseCheckout] Unmarshal Stripe config error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal error: %s", err.Error())
	}

	// Initialize Stripe client with API credentials
	client := stripe.NewClient(stripe.Config{
		SecretKey:     stripeConfig.SecretKey,
		PublicKey:     stripeConfig.PublicKey,
		WebhookSecret: stripeConfig.WebhookSecret,
	})

	// Convert order amount to CNY using current exchange rate
	amount, err := l.queryExchangeRate("CNY", info.Amount)
	if err != nil {
		l.Errorw("[PurchaseCheckout] queryExchangeRate error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "queryExchangeRate error: %s", err.Error())
	}
	convertAmount := int64(amount * 100) // Convert to cents for Stripe API

	// Create Stripe payment sheet for client-side processing
	result, err := client.CreatePaymentSheet(&stripe.Order{
		OrderNo:   info.OrderNo,
		Subscribe: strconv.FormatInt(info.SubscribeId, 10),
		Amount:    convertAmount,
		Currency:  "cny",
		Payment:   stripeConfig.Payment,
	},
		&stripe.User{
			Email: identifier,
		})
	if err != nil {
		l.Errorw("[PurchaseCheckout] CreatePaymentSheet error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "CreatePaymentSheet error: %s", err.Error())
	}

	// Prepare response data for client-side Stripe integration
	stripePayment := &types.StripePayment{
		PublishableKey: stripeConfig.PublicKey,
		ClientSecret:   result.ClientSecret,
		Method:         stripeConfig.Payment,
	}

	// Save Stripe trade number to order for tracking
	info.TradeNo = result.TradeNo
	err = l.deps.OrderModel.Update(l.ctx, info)
	if err != nil {
		l.Errorw("[PurchaseCheckout] Update order error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Update error: %s", err.Error())
	}
	return stripePayment, nil
}

// epayPayment processes EPay payment by generating a payment URL for redirect
// It handles currency conversion and creates a payment URL for external payment processing
func (l *PurchaseCheckoutLogic) epayPayment(paymentConfig *payment.Payment, info *order.Order, returnUrl string) (string, error) {
	var err error
	// Parse EPay configuration from payment settings
	epayConfig := &payment.EPayConfig{}
	if err := epayConfig.Unmarshal([]byte(paymentConfig.Config)); err != nil {
		l.Errorw("[PurchaseCheckout] Unmarshal EPay config error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal error: %s", err.Error())
	}
	// Initialize EPay client with merchant credentials
	client := epay.NewClient(epayConfig.Pid, epayConfig.Url, epayConfig.Key, epayConfig.Type)
	var amount float64
	if l.deps.Config.Currency.Unit != "CNY" {
		// Convert order amount to CNY using current exchange rate
		amount, err = l.queryExchangeRate("CNY", info.Amount)
		if err != nil {
			l.Error("[PurchaseCheckout] queryExchangeRate error", logger.Field("error", err.Error()))
			return "", err
		}
	} else {
		amount = float64(info.Amount) / float64(100)
	}

	// gateway mod
	isGatewayMod := report.IsGatewayMode()

	// Build notification URL for payment status callbacks
	notifyUrl := ""
	if paymentConfig.Domain != "" {
		notifyUrl = paymentConfig.Domain
		if isGatewayMod {
			notifyUrl += "/api/"
		}
		notifyUrl = notifyUrl + "/api/v1/notify/" + paymentConfig.Platform + "/" + paymentConfig.Token
	} else {
		host, ok := l.ctx.Value(config.CtxKeyRequestHost).(string)
		if !ok {
			host = l.deps.Config.Host
		}
		notifyUrl = "https://" + host
		if isGatewayMod {
			notifyUrl += "/api"
		}
		notifyUrl = notifyUrl + "/api/v1/notify/" + paymentConfig.Platform + "/" + paymentConfig.Token
	}

	// Create payment URL for user redirection
	url := client.CreatePayUrl(epay.Order{
		Name:      l.deps.Config.Site.SiteName,
		Amount:    amount,
		OrderNo:   info.OrderNo,
		SignType:  "MD5",
		NotifyUrl: notifyUrl,
		ReturnUrl: returnUrl,
	})
	return url, nil
}

// CryptoSaaSPayment processes CryptoSaaSPayment payment by generating a payment URL for redirect
// It handles currency conversion and creates a payment URL for external payment processing
func (l *PurchaseCheckoutLogic) CryptoSaaSPayment(paymentConfig *payment.Payment, info *order.Order, returnUrl string) (string, error) {
	var err error
	// Parse EPay configuration from payment settings
	epayConfig := &payment.CryptoSaaSConfig{}
	if err := epayConfig.Unmarshal([]byte(paymentConfig.Config)); err != nil {
		l.Errorw("[PurchaseCheckout] Unmarshal EPay config error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unmarshal error: %s", err.Error())
	}
	// Initialize EPay client with merchant credentials
	client := epay.NewClient(epayConfig.AccountID, epayConfig.Endpoint, epayConfig.SecretKey, epayConfig.Type)

	var amount float64

	if l.deps.Config.Currency.Unit != "CNY" {
		// Convert order amount to CNY using current exchange rate
		amount, err = l.queryExchangeRate("CNY", info.Amount)
		if err != nil {
			return "", err
		}
	} else {
		amount = float64(info.Amount) / float64(100)
	}

	// gateway mod
	isGatewayMod := report.IsGatewayMode()

	// Build notification URL for payment status callbacks
	notifyUrl := ""
	if paymentConfig.Domain != "" {
		notifyUrl = paymentConfig.Domain
		if isGatewayMod {
			notifyUrl += "/api/"
		}
		notifyUrl = notifyUrl + "/api/v1/notify/" + paymentConfig.Platform + "/" + paymentConfig.Token
	} else {
		host, ok := l.ctx.Value(config.CtxKeyRequestHost).(string)
		if !ok {
			host = l.deps.Config.Host
		}

		notifyUrl = "https://" + host
		if isGatewayMod {
			notifyUrl += "/api"
		}
		notifyUrl = notifyUrl + "/api/v1/notify/" + paymentConfig.Platform + "/" + paymentConfig.Token
	}
	// Create payment URL for user redirection
	url := client.CreatePayUrl(epay.Order{
		Name:      l.deps.Config.Site.SiteName,
		Amount:    amount,
		OrderNo:   info.OrderNo,
		SignType:  "MD5",
		NotifyUrl: notifyUrl,
		ReturnUrl: returnUrl,
	})
	return url, nil
}

// queryExchangeRate converts the order amount from system currency to target currency
// It retrieves the current exchange rate and performs currency conversion if needed
func (l *PurchaseCheckoutLogic) queryExchangeRate(to string, src int64) (amount float64, err error) {
	// Convert cents to decimal amount
	amount = float64(src) / float64(100)

	// No conversion needed if target currency matches system currency
	if to == l.deps.Config.Currency.Unit {
		return amount, nil
	}

	currentUnit := l.deps.Config.Currency.Unit
	snapshot := l.deps.CurrentExchangeRateSnapshot()
	if snapshot.Rate != 0 && snapshot.From == currentUnit && snapshot.To == to {
		amount = amount * snapshot.Rate
		return amount, nil
	}

	// Skip conversion if no exchange rate API key configured
	if l.deps.Config.Currency.AccessKey == "" {
		return amount, nil
	}

	version := l.deps.PrepareExchangeRateCache(currentUnit, to)

	// Convert currency if system currency differs from target currency
	result, err := exchangeRate.GetExchangeRete(currentUnit, to, l.deps.Config.Currency.AccessKey, 1)
	if err != nil {
		l.Error("[PurchaseCheckout] QueryExchangeRate error", logger.Field("error", err.Error()))
		return 0, err
	}
	if !l.deps.StoreExchangeRateCache(version, currentUnit, to, result) {
		l.Debug("[PurchaseCheckout] Skip stale exchange rate cache write",
			logger.Field("currency", currentUnit),
			logger.Field("target", to),
			logger.Field("version", version),
		)
	}
	return result * amount, nil
}

// balancePayment processes balance payment with gift amount priority logic
// It prioritizes using gift amount first, then regular balance, and creates proper audit logs
func (l *PurchaseCheckoutLogic) balancePayment(u *user.User, o *order.Order) error {
	var userInfo user.User
	var err error
	if o.Amount == 0 {
		// No payment required for zero-amount orders
		l.Info(
			"[PurchaseCheckout] No payment required for zero-amount order",
			logger.Field("orderNo", o.OrderNo),
			logger.Field("userId", u.Id),
		)
		err = l.deps.OrderModel.UpdateOrderStatus(l.ctx, o.OrderNo, 2)
		if err != nil {
			l.Errorw("[PurchaseCheckout] Update order status error",
				logger.Field("error", err.Error()),
				logger.Field("orderNo", o.OrderNo),
				logger.Field("userId", u.Id))
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Update order status error: %s", err.Error())
		}
		goto activation
	}

	err = l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		// Retrieve latest user information with row-level locking
		err := db.Model(&user.User{}).Where("id = ?", u.Id).First(&userInfo).Error
		if err != nil {
			return err
		}

		// Check if user has sufficient total balance (regular + gift)
		totalAvailable := userInfo.Balance + userInfo.GiftAmount
		if totalAvailable < o.Amount {
			return errors.Wrapf(xerr.NewErrCode(xerr.InsufficientBalance),
				"Insufficient balance: required %d, available %d", o.Amount, totalAvailable)
		}

		// Calculate payment distribution: prioritize gift amount first
		var giftUsed, balanceUsed int64
		remainingAmount := o.Amount

		if userInfo.GiftAmount >= remainingAmount {
			// Gift amount covers the entire payment
			giftUsed = remainingAmount
			balanceUsed = 0
		} else {
			// Use all available gift amount, then regular balance
			giftUsed = userInfo.GiftAmount
			balanceUsed = remainingAmount - giftUsed
		}

		// Update user balances
		userInfo.GiftAmount -= giftUsed
		userInfo.Balance -= balanceUsed

		// Save updated user information
		err = l.deps.UserModel.Update(l.ctx, &userInfo)
		if err != nil {
			return err
		}

		// Create gift amount log if gift amount was used
		if giftUsed > 0 {
			giftLog := &log.Gift{
				OrderNo: o.OrderNo,
				Type:    log.GiftTypeReduce, // Type 2 represents gift amount decrease/usage
				Amount:  giftUsed,
				Balance: userInfo.GiftAmount,
				Remark:  "Purchase payment",
			}
			content, _ := giftLog.Marshal()

			err = db.Create(&log.SystemLog{
				Type:     log.TypeGift.Uint8(),
				ObjectID: userInfo.Id,
				Date:     time.Now().Format(time.DateOnly),
				Content:  string(content),
			}).Error
			if err != nil {
				return err
			}
		}

		// Create balance log if regular balance was used
		if balanceUsed > 0 {
			balanceLog := &log.Balance{
				Amount:    balanceUsed,
				Type:      log.BalanceTypePayment, // Type 3 represents payment deduction
				OrderNo:   o.OrderNo,
				Balance:   userInfo.Balance,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := balanceLog.Marshal()
			err = db.Create(&log.SystemLog{
				Type:     log.TypeBalance.Uint8(),
				ObjectID: userInfo.Id,
				Date:     time.Now().Format(time.DateOnly),
				Content:  string(content),
			}).Error
			if err != nil {
				return err
			}
		}

		// Store gift amount used in order for potential refund tracking
		o.GiftAmount = giftUsed
		err = l.deps.OrderModel.Update(l.ctx, o, db)
		if err != nil {
			return err
		}

		// Mark order as paid (status = 2)
		return l.deps.OrderModel.UpdateOrderStatus(l.ctx, o.OrderNo, 2, db)
	})

	if err != nil {
		l.Errorw("[PurchaseCheckout] Balance payment transaction error",
			logger.Field("error", err.Error()),
			logger.Field("orderNo", o.OrderNo),
			logger.Field("userId", u.Id))
		return err
	}

activation:
	// Enqueue order activation task for immediate processing
	payload := queueType.ForthwithActivateOrderPayload{
		OrderNo: o.OrderNo,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		l.Errorw("[PurchaseCheckout] Marshal activation payload error", logger.Field("error", err.Error()))
		return err
	}

	task := asynq.NewTask(queueType.ForthwithActivateOrder, bytes)
	_, err = l.deps.Queue.EnqueueContext(l.ctx, task)
	if err != nil {
		l.Errorw("[PurchaseCheckout] Enqueue activation task error", logger.Field("error", err.Error()))
		return err
	}

	l.Info("[PurchaseCheckout] Balance payment completed successfully",
		logger.Field("orderNo", o.OrderNo),
		logger.Field("userId", u.Id))
	return nil
}
