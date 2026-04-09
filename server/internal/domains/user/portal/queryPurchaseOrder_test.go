package portal

import (
	"context"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	serverconfig "github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	billingrepo "github.com/perfect-panel/server/internal/platform/persistence/billing"
	modelorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	modelpayment "github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

func TestQueryPurchaseOrderWorksWithNormalizedBillingSchema(t *testing.T) {
	db := openPortalBillingCompatDB(t)
	rds := openPortalBillingCompatRedis(t)

	require.NoError(t, db.Create(&modelsubscribe.Subscribe{
		Id:          101,
		Name:        "Normalized Starter",
		Language:    "zh-CN",
		UnitTime:    "month",
		SpeedLimit:  256,
		DeviceLimit: 4,
	}).Error)

	enabled := true
	require.NoError(t, db.Create(&billingrepo.PaymentGateway{
		ID:           201,
		Name:         "Stripe",
		Platform:     "Stripe",
		Token:        "gateway_tok_1",
		Enable:       &enabled,
		PublicConfig: `{"public_key":"pk_live"}`,
	}).Error)
	require.NoError(t, db.Create(&billingrepo.Order{
		ID:               301,
		OrderNo:          "normalized-order",
		UserID:           401,
		Type:             1,
		Quantity:         1,
		Price:            1200,
		Amount:           1080,
		Discount:         120,
		PaymentGatewayID: 201,
		Method:           "Stripe",
		Status:           1,
		SubscribeID:      101,
		IsNew:            true,
	}).Error)
	require.NoError(t, db.Create(&billingrepo.OrderItem{
		OrderID:     301,
		SubscribeID: 101,
		Quantity:    1,
		UnitPrice:   1200,
		Amount:      1080,
	}).Error)

	logic := NewQueryPurchaseOrderLogic(context.Background(), Deps{
		OrderModel:     modelorder.NewModel(db, rds),
		PaymentModel:   modelpayment.NewModel(db, rds),
		SubscribeModel: modelsubscribe.NewModel(db, rds),
		Config: &serverconfig.Config{
			JwtAuth: serverconfig.JwtAuth{
				AccessSecret: "test-secret",
				AccessExpire: 60,
			},
		},
		Redis: rds,
	})

	resp, err := logic.QueryPurchaseOrder(&types.QueryPurchaseOrderRequest{OrderNo: "normalized-order"})

	require.NoError(t, err)
	require.Equal(t, "normalized-order", resp.OrderNo)
	require.EqualValues(t, 1080, resp.Amount)
	require.EqualValues(t, 120, resp.Discount)
	require.Equal(t, "Normalized Starter", resp.Subscribe.Name)
	require.Equal(t, "Stripe", resp.Payment.Platform)
	require.Equal(t, "Stripe", resp.Payment.Name)
}

func openPortalBillingCompatDB(t *testing.T) *gorm.DB {
	t.Helper()

	schemarevisions.RegisterEmbedded()
	db, err := gorm.Open(sqliteDriver.Open("file:"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, schema.Bootstrap(db, schema.DefaultRevisionSource))
	require.NoError(t, db.AutoMigrate(&modelsubscribe.Subscribe{}))
	require.NoError(t, schema.ApplyRevisions(db, schema.DefaultRevisionSource))
	return db
}

func openPortalBillingCompatRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}
