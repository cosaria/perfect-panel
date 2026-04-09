package revisions

import (
	"reflect"
	"strings"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
)

func TestAsyncTrustAndUsageRevisionAvoidsMySQLIncompatibleDefaults(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		model  any
		fields []string
	}{
		{name: "system.ExternalTrustEvent", model: system.ExternalTrustEvent{}, fields: []string{"FailureReason", "RawPayload"}},
		{name: "node.NodeUsageReport", model: node.NodeUsageReport{}, fields: []string{"RawPayload"}},
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
