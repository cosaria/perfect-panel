package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/perfect-panel/server-v2/internal/app/routing"
	"github.com/perfect-panel/server-v2/internal/domains/access"
)

func TestPublicSessionsContract(t *testing.T) {
	doc := loadOpenAPIBundle(t)
	operation := mustOperation(t, doc, "/api/v1/public/sessions", "post")

	assertOperationHasResponses(t, operation, "201", "400", "401", "422")
	assertNoSecurityRequirement(t, operation)
	assertOperationID(t, operation, "publicAuthSessionCreate")
	assertSuccessEnvelope(t, operation, "201")
}

func TestPasswordResetContract(t *testing.T) {
	doc := loadOpenAPIBundle(t)
	requestOperation := mustOperation(t, doc, "/api/v1/public/password-reset-requests", "post")
	resetOperation := mustOperation(t, doc, "/api/v1/public/password-resets", "post")

	assertOperationHasResponses(t, requestOperation, "201", "400", "422")
	assertOperationHasResponses(t, resetOperation, "200", "400", "409", "422")
	assertOperationDoesNotExposeStatus(t, requestOperation, "202")
	assertOperationDoesNotExposeStatus(t, resetOperation, "202")
	assertOperationID(t, requestOperation, "publicAuthPasswordResetRequestCreate")
	assertOperationID(t, resetOperation, "publicAuthPasswordResetCreate")
	assertSuccessEnvelope(t, resetOperation, "200")
}

func TestUserSessionContractMatchesHTTPDesign(t *testing.T) {
	doc := loadOpenAPIBundle(t)
	listOperation := mustOperation(t, doc, "/api/v1/user/me/sessions", "get")
	deleteOperation := mustOperation(t, doc, "/api/v1/user/me/sessions/{sessionId}", "delete")

	assertOperationHasResponses(t, listOperation, "200", "401")
	assertOperationHasResponses(t, deleteOperation, "200", "401")
	assertOperationDoesNotExposeStatus(t, deleteOperation, "204")
	assertOperationID(t, listOperation, "userAuthSessionList")
	assertOperationID(t, deleteOperation, "userAuthSessionDelete")
	assertOperationUsesSessionAuth(t, listOperation)
	assertOperationUsesSessionAuth(t, deleteOperation)
	assertPathParameter(t, deleteOperation, "sessionId")
	assertSuccessEnvelope(t, deleteOperation, "200")
}

func TestAdminRoutingRequiresSessionBeforePermissions(t *testing.T) {
	mux := http.NewServeMux()
	order := make([]string, 0, 4)

	requireSession := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "session")
			next.ServeHTTP(w, r)
		})
	}
	requirePermissions := func(required ...string) func(http.Handler) http.Handler {
		requiredCopy := append([]string(nil), required...)
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "permissions:"+strings.Join(requiredCopy, ","))
				next.ServeHTTP(w, r)
			})
		}
	}

	routing.RegisterAdmin(mux, access.NewAdminHTTPHandler(), requireSession, requirePermissions)

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	getResp := httptest.NewRecorder()
	mux.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/admin/users 应返回 200，实际 %d", getResp.Code)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", strings.NewReader(`{"email":"user@example.com","status":"active"}`))
	postReq.Header.Set("Content-Type", "application/json")
	postResp := httptest.NewRecorder()
	mux.ServeHTTP(postResp, postReq)
	if postResp.Code != http.StatusCreated {
		t.Fatalf("POST /api/v1/admin/users 应返回 201，实际 %d", postResp.Code)
	}

	wantOrder := []string{
		"session",
		"permissions:" + access.PermissionAdminUsersRead,
		"session",
		"permissions:" + access.PermissionAdminUsersWrite,
	}
	if !reflect.DeepEqual(order, wantOrder) {
		t.Fatalf("admin 路由应先过 RequireSession 再过 RequirePermissions，want=%v got=%v", wantOrder, order)
	}
}

func loadOpenAPIBundle(t *testing.T) map[string]any {
	t.Helper()

	moduleRoot := getModuleRoot(t)
	bundlePath := filepath.Join(moduleRoot, "openapi", "dist", "openapi.json")
	runContractPipeline(t, moduleRoot)

	raw, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatalf("读取 bundle 失败: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("解析 bundle 失败: %v", err)
	}
	return doc
}

func mustOperation(t *testing.T, doc map[string]any, path string, method string) map[string]any {
	t.Helper()

	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatal("bundle 缺少 paths")
	}
	pathItem, ok := paths[path].(map[string]any)
	if !ok {
		t.Fatalf("缺少路径 %s", path)
	}
	operation, ok := pathItem[method].(map[string]any)
	if !ok {
		t.Fatalf("缺少操作 %s %s", method, path)
	}
	return operation
}

func assertOperationHasResponses(t *testing.T, operation map[string]any, statuses ...string) {
	t.Helper()

	responses, ok := operation["responses"].(map[string]any)
	if !ok {
		t.Fatal("operation 缺少 responses")
	}
	for _, status := range statuses {
		if _, ok := responses[status]; !ok {
			t.Fatalf("operation 缺少响应码 %s", status)
		}
	}
}

func assertOperationDoesNotExposeStatus(t *testing.T, operation map[string]any, status string) {
	t.Helper()

	responses, ok := operation["responses"].(map[string]any)
	if !ok {
		t.Fatal("operation 缺少 responses")
	}
	if _, ok := responses[status]; ok {
		t.Fatalf("operation 不应暴露响应码 %s", status)
	}
}

func assertNoSecurityRequirement(t *testing.T, operation map[string]any) {
	t.Helper()

	security, ok := operation["security"].([]any)
	if !ok {
		t.Fatal("operation 缺少 security")
	}
	if len(security) != 0 {
		t.Fatalf("public operation 不应要求安全方案: %+v", security)
	}
}

func assertOperationUsesSessionAuth(t *testing.T, operation map[string]any) {
	t.Helper()

	security, ok := operation["security"].([]any)
	if !ok || len(security) == 0 {
		t.Fatal("operation 缺少 security")
	}

	first, ok := security[0].(map[string]any)
	if !ok {
		t.Fatalf("security 项格式不正确: %+v", security[0])
	}
	if _, ok := first["sessionAuth"]; !ok {
		t.Fatalf("operation 未使用 sessionAuth: %+v", first)
	}
}

func assertPathParameter(t *testing.T, operation map[string]any, name string) {
	t.Helper()

	rawParameters, ok := operation["parameters"].([]any)
	if !ok {
		t.Fatalf("operation 缺少 parameters: %+v", operation)
	}
	for _, raw := range rawParameters {
		parameter, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if parameter["name"] == name && parameter["in"] == "path" {
			return
		}
	}
	t.Fatalf("operation 缺少路径参数 %s", name)
}

func assertSuccessEnvelope(t *testing.T, operation map[string]any, status string) {
	t.Helper()

	responses := operation["responses"].(map[string]any)
	response, ok := responses[status].(map[string]any)
	if !ok {
		t.Fatalf("缺少响应 %s", status)
	}
	content := response["content"].(map[string]any)
	applicationJSON := content["application/json"].(map[string]any)
	schema := applicationJSON["schema"].(map[string]any)
	required := schema["required"].([]any)

	if len(required) != 2 || required[0] != "data" || required[1] != "meta" {
		t.Fatalf("成功响应必须显式包含 data + meta，实际 %+v", required)
	}
}

func assertOperationID(t *testing.T, operation map[string]any, want string) {
	t.Helper()

	got, ok := operation["operationId"].(string)
	if !ok || got == "" {
		t.Fatal("operation 必须显式声明 operationId")
	}
	if got != want {
		t.Fatalf("operationId 不匹配: want=%s got=%s", want, got)
	}
}
