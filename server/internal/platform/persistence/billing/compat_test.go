package billing_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	modelorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	modelpayment "github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	schemarevisions "github.com/perfect-panel/server/internal/platform/persistence/schema/revisions"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modelsubscription "github.com/perfect-panel/server/internal/platform/persistence/subscription"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestOrderAndPaymentCompatibilityPreferNormalizedBillingSchema(t *testing.T) {
	t.Parallel()

	db := openBillingCompatDB(t)
	rds := openBillingCompatRedis(t)
	ctx := context.Background()

	if err := db.Create(&modelsubscribe.Subscribe{
		Id:          11,
		Name:        "Starter",
		Language:    "zh-CN",
		UnitTime:    "month",
		SpeedLimit:  128,
		DeviceLimit: 3,
	}).Error; err != nil {
		t.Fatalf("create legacy subscribe catalog: %v", err)
	}

	enabled := true
	if err := db.Create(&billing.PaymentGateway{
		ID:           21,
		Name:         "Stripe",
		Platform:     "Stripe",
		Token:        "pay_tok_1",
		Enable:       &enabled,
		PublicConfig: `{"public_key":"pk_test"}`,
	}).Error; err != nil {
		t.Fatalf("create normalized payment gateway: %v", err)
	}
	if err := db.Create(&billing.PaymentGatewaySecret{
		PaymentGatewayID: 21,
		SecretConfig:     `{"secret_key":"sk_test"}`,
	}).Error; err != nil {
		t.Fatalf("create normalized payment gateway secret: %v", err)
	}
	if err := db.Create(&billing.Order{
		ID:               31,
		OrderNo:          "order_norm_1",
		UserID:           41,
		Type:             1,
		Quantity:         1,
		Price:            1000,
		Amount:           900,
		Discount:         100,
		PaymentGatewayID: 21,
		Method:           "Stripe",
		Status:           1,
		SubscribeID:      11,
		IsNew:            true,
	}).Error; err != nil {
		t.Fatalf("create normalized order: %v", err)
	}
	if err := db.Create(&billing.OrderItem{
		OrderID:     31,
		SubscribeID: 11,
		Quantity:    1,
		UnitPrice:   1000,
		Amount:      900,
	}).Error; err != nil {
		t.Fatalf("create normalized order item: %v", err)
	}

	paymentModel := modelpayment.NewModel(db, rds)
	methods, err := paymentModel.FindAvailableMethods(ctx)
	if err != nil {
		t.Fatalf("FindAvailableMethods returned error: %v", err)
	}
	if len(methods) != 1 || methods[0].Token != "pay_tok_1" {
		t.Fatalf("expected normalized payment gateway to surface through compatibility model, got %+v", methods)
	}

	method, err := paymentModel.FindOneByPaymentToken(ctx, "pay_tok_1")
	if err != nil {
		t.Fatalf("FindOneByPaymentToken returned error: %v", err)
	}
	if method.Platform != "Stripe" {
		t.Fatalf("expected Stripe platform, got %+v", method.Platform)
	}

	orderModel := modelorder.NewModel(db, rds)
	total, list, err := orderModel.QueryOrderListByPage(ctx, 1, 10, 0, 0, 0, "order_norm_1")
	if err != nil {
		t.Fatalf("QueryOrderListByPage returned error: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected one normalized order from compatibility facade, got total=%d list=%d", total, len(list))
	}
	if list[0].Payment == nil || list[0].Payment.Token != "pay_tok_1" {
		t.Fatalf("expected payment preload from normalized schema, got %+v", list[0].Payment)
	}
	if list[0].Subscribe == nil || list[0].Subscribe.Id != 11 {
		t.Fatalf("expected subscribe preload from catalog facade, got %+v", list[0].Subscribe)
	}
}

func TestUserSubscriptionCompatibilityReadsNormalizedSubscriptionSchema(t *testing.T) {
	t.Parallel()

	db := openBillingCompatDB(t)
	rds := openBillingCompatRedis(t)
	ctx := context.Background()

	if err := db.Create(&modelsubscription.Subscription{
		ID:          51,
		UserID:      61,
		OrderID:     71,
		SubscribeID: 81,
		Status:      1,
		Traffic:     2048,
		Download:    512,
		Upload:      256,
	}).Error; err != nil {
		t.Fatalf("create normalized subscription: %v", err)
	}
	if err := db.Create(&modelsubscription.SubscriptionToken{
		SubscriptionID: 51,
		Token:          "sub_tok_1",
		UUID:           "uuid-sub-1",
		IsPrimary:      true,
	}).Error; err != nil {
		t.Fatalf("create normalized subscription token: %v", err)
	}

	userModel := modeluser.NewModel(db, rds)
	userSubscribe, err := userModel.FindOneSubscribeByToken(ctx, "sub_tok_1")
	if err != nil {
		t.Fatalf("FindOneSubscribeByToken returned error: %v", err)
	}
	if userSubscribe.Id != 51 || userSubscribe.UUID != "uuid-sub-1" {
		t.Fatalf("expected normalized subscription token lookup result, got %+v", userSubscribe)
	}

	byOrder, err := userModel.FindOneSubscribeByOrderId(ctx, 71)
	if err != nil {
		t.Fatalf("FindOneSubscribeByOrderId returned error: %v", err)
	}
	if byOrder.Token != "sub_tok_1" {
		t.Fatalf("expected order lookup to hydrate token from normalized schema, got %+v", byOrder.Token)
	}
}

func openBillingCompatDB(t *testing.T) *gorm.DB {
	t.Helper()

	schemarevisions.RegisterEmbedded()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()),
	)
	db, err := gorm.Open(sqliteDriver.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := schema.Bootstrap(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("bootstrap schema: %v", err)
	}
	if err := db.AutoMigrate(&modelsubscribe.Subscribe{}); err != nil {
		t.Fatalf("migrate subscribe catalog table: %v", err)
	}
	if err := schema.ApplyRevisions(db, schema.DefaultRevisionSource); err != nil {
		t.Fatalf("apply revisions: %v", err)
	}
	return db
}

func openBillingCompatRedis(t *testing.T) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}
