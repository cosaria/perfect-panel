package portal

import (
	"context"
	"errors"
	"testing"

	modelorder "github.com/perfect-panel/server/models/order"
	modelpayment "github.com/perfect-panel/server/models/payment"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/stretchr/testify/require"
)

func TestPurchaseCheckoutReturnsOrderNotExistWhenLookupFails(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return nil, errors.New("lookup failed")
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "missing"})

	requirePortalCodeError(t, err, xerr.OrderNotExist)
}

func TestPurchaseCheckoutRejectsNonPendingOrder(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo: "paid",
				Status:  2,
			}, nil
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "paid"})

	requirePortalCodeError(t, err, xerr.OrderStatusError)
}

func TestPurchaseCheckoutReturnsDatabaseQueryErrorWhenPaymentLookupFails(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo:   "pending",
				Status:    1,
				PaymentId: 9,
				Method:    "Stripe",
			}, nil
		},
	}
	deps.PaymentModel = fakePortalPaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return nil, errors.New("payment lookup failed")
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "pending"})

	requirePortalCodeError(t, err, xerr.DatabaseQueryError)
}

func TestPurchaseCheckoutRoutesStripeBranchOnConfigParseError(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo:   "stripe-order",
				Status:    1,
				PaymentId: 1,
				Method:    "Stripe",
			}, nil
		},
	}
	deps.PaymentModel = fakePortalPaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return &modelpayment.Payment{
				Id:       1,
				Platform: "Stripe",
				Config:   "{invalid-json",
			}, nil
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "stripe-order"})

	requirePortalCodeError(t, err, xerr.ERROR)
	require.ErrorContains(t, err, "stripePayment error: Unmarshal error")
}

func TestPurchaseCheckoutRoutesEPayBranchOnConfigParseError(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo:   "epay-order",
				Status:    1,
				PaymentId: 2,
				Method:    "EPay",
			}, nil
		},
	}
	deps.PaymentModel = fakePortalPaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return &modelpayment.Payment{
				Id:       2,
				Platform: "EPay",
				Config:   "{invalid-json",
			}, nil
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "epay-order"})

	requirePortalCodeError(t, err, xerr.ERROR)
	require.ErrorContains(t, err, "epayPayment error: Unmarshal error")
}

func TestPurchaseCheckoutRoutesAlipayF2FBranchOnConfigParseError(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo:   "alipay-order",
				Status:    1,
				PaymentId: 3,
				Method:    "AlipayF2F",
			}, nil
		},
	}
	deps.PaymentModel = fakePortalPaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return &modelpayment.Payment{
				Id:       3,
				Platform: "AlipayF2F",
				Config:   "{invalid-json",
			}, nil
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "alipay-order"})

	requirePortalCodeError(t, err, xerr.ERROR)
	require.ErrorContains(t, err, "alipayF2fPayment error: Unmarshal error")
}

func TestPurchaseCheckoutRejectsBalanceCheckoutWithoutUserID(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo:   "balance-order",
				Status:    1,
				PaymentId: 4,
				Method:    "balance",
				UserId:    0,
			}, nil
		},
	}
	deps.PaymentModel = fakePortalPaymentModel{
		findOneFn: func(context.Context, int64) (*modelpayment.Payment, error) {
			return &modelpayment.Payment{
				Id:       4,
				Platform: "balance",
				Config:   "{}",
			}, nil
		},
	}
	logic := NewPurchaseCheckoutLogic(context.Background(), deps)

	_, err := logic.PurchaseCheckout(&types.CheckoutOrderRequest{OrderNo: "balance-order"})

	requirePortalCodeError(t, err, xerr.UserNotExist)
}
