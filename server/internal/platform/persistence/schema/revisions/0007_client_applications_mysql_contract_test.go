package revisions

import (
	"reflect"
	"strings"
	"testing"

	"github.com/perfect-panel/server/internal/platform/persistence/client"
)

func TestClientApplicationsRevisionAvoidsMySQLIncompatibleDefaults(t *testing.T) {
	t.Parallel()

	cases := []struct {
		fieldName string
		expect    string
	}{
		{fieldName: "Icon", expect: "type:MEDIUMTEXT"},
		{fieldName: "SubscribeTemplate", expect: "type:MEDIUMTEXT"},
		{fieldName: "DownloadLink", expect: "type:text"},
	}

	modelType := reflect.TypeOf(client.SubscribeApplication{})
	for _, tc := range cases {
		field, ok := modelType.FieldByName(tc.fieldName)
		if !ok {
			t.Fatalf("SubscribeApplication 缺少字段 %s", tc.fieldName)
		}
		tag := field.Tag.Get("gorm")
		if !strings.Contains(tag, tc.expect) {
			t.Fatalf("%s 预期包含 %q，实际 tag=%q", tc.fieldName, tc.expect, tag)
		}
		if strings.Contains(tag, "default:") {
			t.Fatalf("%s 不应为 text/mediumtext 列声明默认值，实际 tag=%q", tc.fieldName, tag)
		}
	}
}
