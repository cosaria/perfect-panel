package order

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	modelcoupon "github.com/perfect-panel/server/models/coupon"
	modelpayment "github.com/perfect-panel/server/models/payment"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPurchaseReturnsInvalidAccessWhenUserMissingFromContext(t *testing.T) {
	logic := NewPurchaseLogic(context.Background(), newOrderTestDeps())

	_, err := logic.Purchase(&types.PurchaseOrderRequest{Quantity: 1})

	requireCodeError(t, err, xerr.InvalidAccess)
}

func TestPurchaseNormalizesQuantityToOneBeforeFurtherValidation(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		SubscribeId: 101,
		Quantity:    0,
		Payment:     9,
	}
	deps := newOrderTestDeps()
	deps.UserModel = fakeUserModel{
		queryUserSubscribeFn: func(context.Context, int64, ...int64) ([]*modeluser.SubscribeDetails, error) {
			return nil, nil
		},
	}
	deps.SubscribeModel = fakeSubscribeModel{
		findOneFn: func(context.Context, int64) (*modelsubscribe.Subscribe, error) {
			return &modelsubscribe.Subscribe{
				Id:        req.SubscribeId,
				UnitPrice: 1000,
				Inventory: 10,
				Sell:      boolPtr(true),
			}, nil
		},
	}
	deps.PaymentModel = fakePaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return nil, errors.New("payment lookup failed")
		},
	}

	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 7})
	logic := NewPurchaseLogic(ctx, deps)

	_, err := logic.Purchase(req)

	requireCodeError(t, err, xerr.DatabaseQueryError)
	require.Equal(t, int64(1), req.Quantity)
}

func TestPurchaseRejectsQuantityAboveMax(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		Quantity: MaxQuantity + 1,
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 8})
	logic := NewPurchaseLogic(ctx, newOrderTestDeps())

	_, err := logic.Purchase(req)

	requireCodeError(t, err, xerr.InvalidParams)
	require.ErrorContains(t, err, "quantity exceeds maximum limit")
}

func TestPurchaseRejectsOutOfStockSubscribe(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		SubscribeId: 33,
		Quantity:    1,
	}
	deps := newOrderTestDeps()
	deps.UserModel = fakeUserModel{
		queryUserSubscribeFn: func(context.Context, int64, ...int64) ([]*modeluser.SubscribeDetails, error) {
			return nil, nil
		},
	}
	deps.SubscribeModel = fakeSubscribeModel{
		findOneFn: func(context.Context, int64) (*modelsubscribe.Subscribe, error) {
			return &modelsubscribe.Subscribe{
				Id:        req.SubscribeId,
				UnitPrice: 1000,
				Inventory: 0,
				Sell:      boolPtr(true),
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 9})
	logic := NewPurchaseLogic(ctx, deps)

	_, err := logic.Purchase(req)

	requireCodeError(t, err, xerr.SubscribeOutOfStock)
}

func TestPurchaseRejectsQuotaLimit(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		SubscribeId: 77,
		Quantity:    1,
	}
	deps := newOrderTestDeps()
	deps.UserModel = fakeUserModel{
		queryUserSubscribeFn: func(context.Context, int64, ...int64) ([]*modeluser.SubscribeDetails, error) {
			return []*modeluser.SubscribeDetails{
				{SubscribeId: req.SubscribeId},
			}, nil
		},
	}
	deps.SubscribeModel = fakeSubscribeModel{
		findOneFn: func(context.Context, int64) (*modelsubscribe.Subscribe, error) {
			return &modelsubscribe.Subscribe{
				Id:        req.SubscribeId,
				UnitPrice: 1000,
				Inventory: 10,
				Quota:     1,
				Sell:      boolPtr(true),
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 10})
	logic := NewPurchaseLogic(ctx, deps)

	_, err := logic.Purchase(req)

	requireCodeError(t, err, xerr.SubscribeQuotaLimit)
}

func TestPurchaseReturnsCouponNotExistWhenCouponLookupMisses(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		SubscribeId: 55,
		Quantity:    1,
		Coupon:      "MISSING",
	}
	deps := newOrderTestDeps()
	deps.UserModel = fakeUserModel{
		queryUserSubscribeFn: func(context.Context, int64, ...int64) ([]*modeluser.SubscribeDetails, error) {
			return nil, nil
		},
	}
	deps.SubscribeModel = fakeSubscribeModel{
		findOneFn: func(context.Context, int64) (*modelsubscribe.Subscribe, error) {
			return &modelsubscribe.Subscribe{
				Id:        req.SubscribeId,
				UnitPrice: 1000,
				Inventory: 10,
				Sell:      boolPtr(true),
			}, nil
		},
	}
	deps.CouponModel = fakeCouponModel{
		findOneByCodeFn: func(context.Context, string) (*modelcoupon.Coupon, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 11})
	logic := NewPurchaseLogic(ctx, deps)

	_, err := logic.Purchase(req)

	requireCodeError(t, err, xerr.CouponNotExist)
}
