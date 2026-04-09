package order

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PreCreateOrderInput struct {
	Body types.PurchaseOrderRequest
}

type PreCreateOrderOutput struct {
	Body *types.PreOrderResponse
}

func PreCreateOrderHandler(deps Deps) func(context.Context, *PreCreateOrderInput) (*PreCreateOrderOutput, error) {
	return func(ctx context.Context, input *PreCreateOrderInput) (*PreCreateOrderOutput, error) {
		l := NewPreCreateOrderLogic(ctx, deps)
		resp, err := l.PreCreateOrder(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PreCreateOrderOutput{Body: resp}, nil
	}
}

type PreCreateOrderLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewPreCreateOrderLogic creates a new pre-create order logic instance for order preview operations.
// It initializes the logger with context and sets up the service context for database operations.
func NewPreCreateOrderLogic(ctx context.Context, deps Deps) *PreCreateOrderLogic {
	return &PreCreateOrderLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

// PreCreateOrder calculates order pricing preview including discounts, coupons, gift amounts, and fees
// without actually creating an order. It validates subscription plans, coupons, and payment methods
// to provide accurate pricing information for the frontend order preview.
func (l *PreCreateOrderLogic) PreCreateOrder(req *types.PurchaseOrderRequest) (resp *types.PreOrderResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	if req.Quantity <= 0 {
		l.Debugf("[PreCreateOrder] Quantity is less than or equal to 0, setting to 1")
		req.Quantity = 1
	}

	// find subscribe plan
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, req.SubscribeId)
	if err != nil {
		l.Errorw("[PreCreateOrder] Database query error", logger.Field("error", err.Error()), logger.Field("subscribe_id", req.SubscribeId))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find subscribe error: %v", err.Error())
	}
	var discount float64 = 1
	if sub.Discount != "" {
		var dis []types.SubscribeDiscount
		_ = json.Unmarshal([]byte(sub.Discount), &dis)
		discount = getDiscount(dis, req.Quantity)
	}
	price := sub.UnitPrice * req.Quantity

	amount := int64(float64(price) * discount)
	discountAmount := price - amount
	var couponAmount int64
	if req.Coupon != "" {
		couponInfo, err := l.deps.CouponModel.FindOneByCode(l.ctx, req.Coupon)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.CouponNotExist), "coupon not found")
			}
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find coupon error: %v", err.Error())
		}
		if couponInfo.Count > 0 && couponInfo.Count <= couponInfo.UsedCount {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.CouponAlreadyUsed), "coupon used")
		}
		if l.deps.DB == nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "order db is nil")
		}
		var count int64
		err = l.deps.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Model(&order.Order{}).Where("user_id = ? and coupon = ?", u.Id, req.Coupon).Count(&count).Error
		})

		if err != nil {
			l.Errorw("[PreCreateOrder] Database query error", logger.Field("error", err.Error()), logger.Field("user_id", u.Id), logger.Field("coupon", req.Coupon))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find coupon error: %v", err.Error())
		}

		if couponInfo.UserLimit > 0 && count >= couponInfo.UserLimit {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.CouponInsufficientUsage), "coupon limit exceeded")
		}

		couponSub := tool.StringToInt64Slice(couponInfo.Subscribe)
		if len(couponSub) > 0 && !tool.Contains(couponSub, req.SubscribeId) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.CouponNotApplicable), "coupon not match")
		}
		couponAmount = calculateCoupon(amount, couponInfo)
	}
	amount -= couponAmount

	var deductionAmount int64
	// Check user deduction amount
	if u.GiftAmount > 0 {
		if u.GiftAmount >= amount {
			deductionAmount = amount
			amount = 0
		} else {
			deductionAmount = u.GiftAmount
			amount -= u.GiftAmount
		}
	}
	var feeAmount int64
	if req.Payment != 0 {
		payment, err := l.deps.PaymentModel.FindOne(l.ctx, req.Payment)
		if err != nil {
			l.Errorw("[PreCreateOrder] Database query error", logger.Field("error", err.Error()), logger.Field("payment", req.Payment))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find payment method error: %v", err.Error())
		}
		// Calculate the handling fee
		if amount > 0 {
			feeAmount = calculateFee(amount, payment)
		}
		amount += feeAmount
	}

	resp = &types.PreOrderResponse{
		Price:          price,
		Amount:         amount,
		Discount:       discountAmount,
		GiftAmount:     deductionAmount,
		Coupon:         req.Coupon,
		CouponDiscount: couponAmount,
		FeeAmount:      feeAmount,
	}
	return
}
