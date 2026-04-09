package order

import (
	"context"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	modelcoupon "github.com/perfect-panel/server/internal/platform/persistence/coupon"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
)

func TestPreCreateOrderReturnsErrorInsteadOfPanicWhenCouponPathHasNoDB(t *testing.T) {
	req := &types.PurchaseOrderRequest{
		SubscribeId: 88,
		Quantity:    1,
		Coupon:      "WELCOME",
	}
	deps := newOrderTestDeps()
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
			return &modelcoupon.Coupon{
				Code:      req.Coupon,
				Type:      2,
				Discount:  100,
				UserLimit: 1,
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 12})
	logic := NewPreCreateOrderLogic(ctx, deps)

	require.NotPanics(t, func() {
		_, err := logic.PreCreateOrder(req)
		requireCodeError(t, err, xerr.DatabaseQueryError)
		require.ErrorContains(t, err, "order db is nil")
	})
}
