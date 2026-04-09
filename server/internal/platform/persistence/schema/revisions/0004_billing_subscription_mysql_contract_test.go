package revisions

import (
	"reflect"
	"strings"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	"github.com/perfect-panel/server/internal/platform/persistence/subscription"
)

func TestBillingSubscriptionRevisionAvoidsMySQLIncompatibleDefaults(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		model  any
		fields []string
	}{
		{name: "billing.OrderItem", model: billing.OrderItem{}, fields: []string{"Snapshot"}},
		{name: "billing.PaymentGateway", model: billing.PaymentGateway{}, fields: []string{"PublicConfig", "Description"}},
		{name: "billing.PaymentGatewaySecret", model: billing.PaymentGatewaySecret{}, fields: []string{"SecretConfig"}},
		{name: "billing.Payment", model: billing.Payment{}, fields: []string{"RawPayload"}},
		{name: "billing.PaymentCallback", model: billing.PaymentCallback{}, fields: []string{"RawPayload"}},
		{name: "billing.Refund", model: billing.Refund{}, fields: []string{"Reason"}},
		{name: "subscription.SubscriptionEvent", model: subscription.SubscriptionEvent{}, fields: []string{"Payload"}},
	}

	for _, tc := range cases {
		modelType := reflect.TypeOf(tc.model)
		for _, fieldName := range tc.fields {
			field, ok := modelType.FieldByName(fieldName)
			if !ok {
				t.Fatalf("%s 缺少字段 %s", tc.name, fieldName)
			}
			tag := field.Tag.Get("gorm")
			if !strings.Contains(tag, "type:text") {
				t.Fatalf("%s.%s 预期使用 text 列，实际 tag=%q", tc.name, fieldName, tag)
			}
			if strings.Contains(tag, "default:") {
				t.Fatalf("%s.%s 不应为 text 列声明默认值，实际 tag=%q", tc.name, fieldName, tag)
			}
		}
	}
}

func TestLegacyUserSubscriptionTableUsesMySQLCompatibleStartTime(t *testing.T) {
	t.Parallel()

	field, ok := reflect.TypeOf(legacyUserSubscriptionTable{}).FieldByName("StartTime")
	if !ok {
		t.Fatal("legacyUserSubscriptionTable 缺少 StartTime 字段")
	}

	tag := field.Tag.Get("gorm")
	if !strings.Contains(tag, "type:datetime") {
		t.Fatalf("StartTime 需要显式使用 datetime，实际 tag=%q", tag)
	}
	if !strings.Contains(tag, "default:CURRENT_TIMESTAMP") {
		t.Fatalf("StartTime 需要保留 CURRENT_TIMESTAMP 默认值，实际 tag=%q", tag)
	}
}
