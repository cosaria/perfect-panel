package portal

import (
	"context"
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	modelorder "github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
)

func TestQueryPurchaseOrderReturnsErrorInsteadOfPanicWhenRedisMissingForTemporaryOrder(t *testing.T) {
	deps := newPortalTestDeps()
	deps.OrderModel = fakePortalOrderModel{
		findOneByOrderNoFn: func(context.Context, string) (*modelorder.Order, error) {
			return &modelorder.Order{
				OrderNo: "temp-order",
				Status:  2,
			}, nil
		},
	}
	logic := NewQueryPurchaseOrderLogic(context.Background(), deps)

	require.NotPanics(t, func() {
		_, err := logic.QueryPurchaseOrder(&types.QueryPurchaseOrderRequest{OrderNo: "temp-order"})
		requirePortalCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "redis client is nil")
	})
}

func TestGenerateSessionTokenReturnsErrorInsteadOfPanicWhenRedisMissing(t *testing.T) {
	deps := Deps{
		Config: &serverconfig.Config{
			JwtAuth: serverconfig.JwtAuth{
				AccessSecret: "test-secret",
				AccessExpire: 60,
			},
		},
	}
	logic := NewQueryPurchaseOrderLogic(context.Background(), deps)

	require.NotPanics(t, func() {
		_, err := logic.generateSessionToken(42)
		requirePortalCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "redis client is nil")
	})
}

func TestGenerateSessionTokenReturnsErrorInsteadOfPanicWhenConfigMissing(t *testing.T) {
	logic := NewQueryPurchaseOrderLogic(context.Background(), Deps{})

	require.NotPanics(t, func() {
		_, err := logic.generateSessionToken(42)
		requirePortalCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "config is nil")
	})
}

func TestValidateUserAndEmailReturnsErrorInsteadOfPanicWhenUserModelMissing(t *testing.T) {
	logic := NewQueryPurchaseOrderLogic(context.Background(), Deps{})

	require.NotPanics(t, func() {
		err := logic.validateUserAndEmail(&modelorder.Order{UserId: 42}, "email", "user@example.com")
		requirePortalCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "user model is nil")
	})
}
