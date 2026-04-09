package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/modules/auth/jwt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/modules/util/uuidx"
	"github.com/pkg/errors"
)

type QueryPurchaseOrderInput struct {
	types.QueryPurchaseOrderRequest
}

type QueryPurchaseOrderOutput struct {
	Body *types.QueryPurchaseOrderResponse
}

func QueryPurchaseOrderHandler(deps Deps) func(context.Context, *QueryPurchaseOrderInput) (*QueryPurchaseOrderOutput, error) {
	return func(ctx context.Context, input *QueryPurchaseOrderInput) (*QueryPurchaseOrderOutput, error) {
		l := NewQueryPurchaseOrderLogic(ctx, deps)
		resp, err := l.QueryPurchaseOrder(&input.QueryPurchaseOrderRequest)
		if err != nil {
			return nil, err
		}
		return &QueryPurchaseOrderOutput{Body: resp}, nil
	}
}

type QueryPurchaseOrderLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryPurchaseOrderLogic Query Purchase Order
func NewQueryPurchaseOrderLogic(ctx context.Context, deps Deps) *QueryPurchaseOrderLogic {
	return &QueryPurchaseOrderLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

// Centralized error handler for database issues
func wrapDatabaseError(err error) error {
	return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Database Query Error: %v", err.Error())
}

func (l *QueryPurchaseOrderLogic) QueryPurchaseOrder(req *types.QueryPurchaseOrderRequest) (resp *types.QueryPurchaseOrderResponse, err error) {
	orderInfo, err := l.deps.OrderModel.FindOneByOrderNo(l.ctx, req.OrderNo)
	if err != nil {
		return nil, wrapDatabaseError(err)
	}
	// Handle temporary orders if applicable
	var token string
	if orderInfo.Status == 2 || orderInfo.Status == 5 {
		if token, err = l.handleTemporaryOrder(orderInfo, req); err != nil {
			return nil, err
		}
	}
	// Fetch subscription and payment information
	subscribeInfo, paymentInfo, err := l.fetchOrderDetails(orderInfo)
	if err != nil {
		return nil, err
	}

	return &types.QueryPurchaseOrderResponse{
		OrderNo:        orderInfo.OrderNo,
		Subscribe:      subscribeInfo,
		Quantity:       orderInfo.Quantity,
		Price:          orderInfo.Price,
		Amount:         orderInfo.Amount,
		Discount:       orderInfo.Discount,
		Coupon:         orderInfo.Coupon,
		CouponDiscount: orderInfo.CouponDiscount,
		FeeAmount:      orderInfo.FeeAmount,
		Payment:        paymentInfo,
		Status:         orderInfo.Status,
		CreatedAt:      orderInfo.CreatedAt.UnixMilli(),
		Token:          token,
	}, nil
}

// handleTemporaryOrder processes temporary order-related operations
func (l *QueryPurchaseOrderLogic) handleTemporaryOrder(orderInfo *order.Order, req *types.QueryPurchaseOrderRequest) (string, error) {
	if l.deps.Redis == nil {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "redis client is nil")
	}

	cacheKey := fmt.Sprintf(config.TempOrderCacheKey, orderInfo.OrderNo)
	cacheValue, err := l.deps.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil {
		l.Errorw("Get TempOrderCacheKey Error", logger.Field("cacheKey", cacheKey), logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Get TempOrderCacheKey Error: %v", err.Error())
	}

	var tempOrder config.TemporaryOrderInfo
	if err := json.Unmarshal([]byte(cacheValue), &tempOrder); err != nil {
		l.Errorw("JSON Unmarshal Error", logger.Field("error", err.Error()), logger.Field("cacheValue", cacheValue))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "JSON Unmarshal Error: %v", err.Error())
	}
	if tempOrder.OrderNo != orderInfo.OrderNo {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Order number mismatch")
	}

	// Validate user and email
	if err = l.validateUserAndEmail(orderInfo, req.AuthType, req.Identifier); err != nil {
		return "", err
	}

	// Generate session token
	return l.generateSessionToken(orderInfo.UserId)
}

// validateUserAndEmail ensures the user and email are correct
func (l *QueryPurchaseOrderLogic) validateUserAndEmail(orderInfo *order.Order, platform, openid string) error {
	if l.deps.UserModel == nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "user model is nil")
	}

	userInfo, err := l.deps.UserModel.FindOne(l.ctx, orderInfo.UserId)
	if err != nil {
		return wrapDatabaseError(err)
	}

	authMethod, err := l.deps.UserModel.FindUserAuthMethodByOpenID(l.ctx, platform, openid)
	if err != nil {
		return wrapDatabaseError(err)
	}
	if authMethod.UserId != userInfo.Id {
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Email verification failed")
	}

	return nil
}

// generateSessionToken creates a session token and stores it in Redis
func (l *QueryPurchaseOrderLogic) generateSessionToken(userId int64) (string, error) {
	if l.deps.Config == nil {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "config is nil")
	}

	sessionId := uuidx.NewUUID().String()
	token, err := jwt.NewJwtToken(
		l.deps.Config.JwtAuth.AccessSecret,
		time.Now().Unix(),
		l.deps.Config.JwtAuth.AccessExpire,
		jwt.WithOption("UserId", userId),
		jwt.WithOption("SessionId", sessionId),
	)
	if err != nil {
		l.Errorw("Token Generation Error", logger.Field("error", err.Error()))
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Token generation error")
	}

	if l.deps.Redis == nil {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "redis client is nil")
	}

	cacheKey := fmt.Sprintf("%v:%v", config.SessionIdKey, sessionId)
	if err := l.deps.Redis.Set(l.ctx, cacheKey, userId, time.Duration(l.deps.Config.JwtAuth.AccessExpire)*time.Second).Err(); err != nil {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Session storage error")
	}

	return token, nil
}

// fetchOrderDetails retrieves subscription and payment details
func (l *QueryPurchaseOrderLogic) fetchOrderDetails(orderInfo *order.Order) (types.Subscribe, types.PaymentMethod, error) {
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, orderInfo.SubscribeId)
	if err != nil {
		return types.Subscribe{}, types.PaymentMethod{}, wrapDatabaseError(err)
	}

	var subscribeInfo types.Subscribe
	tool.DeepCopy(&subscribeInfo, sub)

	payment, err := l.deps.PaymentModel.FindOne(l.ctx, orderInfo.PaymentId)
	if err != nil {
		return types.Subscribe{}, types.PaymentMethod{}, wrapDatabaseError(err)
	}

	var paymentInfo types.PaymentMethod
	tool.DeepCopy(&paymentInfo, payment)

	return subscribeInfo, paymentInfo, nil
}
