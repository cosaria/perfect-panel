package contract_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestContractPipelinePasses(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	bundlePath, generatedPath := runContractPipeline(t, moduleRoot)

	if info, err := os.Stat(bundlePath); err != nil {
		t.Fatalf("未生成 bundle 文件 %s: %v", bundlePath, err)
	} else if info.IsDir() {
		t.Fatalf("bundle 路径是目录而不是文件: %s", bundlePath)
	}

	info, err := os.Stat(generatedPath)
	if err != nil {
		t.Fatalf("未生成 openapi-ts 输出目录 %s: %v", generatedPath, err)
	}
	if !info.IsDir() {
		t.Fatalf("openapi-ts 输出路径不是目录: %s", generatedPath)
	}

	entries, err := os.ReadDir(generatedPath)
	if err != nil {
		t.Fatalf("读取 openapi-ts 输出目录失败: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("openapi-ts 输出目录为空: %s", generatedPath)
	}

	assertBundleHasRequiredSurface(t, bundlePath)
	assertNodeUsageReportOperation(t, bundlePath)
	assertPathFilesUseRootComponents(t, moduleRoot)
	assertRootSpecUsesComponentRefs(t, moduleRoot)
	assertGeneratedTypesHaveNoDuplicatePublicAliases(t, generatedPath)
	assertNoOpenApiTsErrorLogs(t, moduleRoot)
}

func getModuleRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法获取当前测试文件路径")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func runContractPipeline(t *testing.T, moduleRoot string) (string, string) {
	t.Helper()

	cleanupGeneratedArtifacts(t, moduleRoot)

	cmd := exec.Command("make", "contract")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make contract 失败: %v\n输出:\n%s", err, string(out))
	}

	t.Cleanup(func() {
		cleanupGeneratedArtifacts(t, moduleRoot)
	})

	return filepath.Join(moduleRoot, "openapi", "dist", "openapi.json"), filepath.Join(moduleRoot, "openapi", "generated", "ts")
}

func cleanupGeneratedArtifacts(t *testing.T, moduleRoot string) {
	t.Helper()

	for _, path := range []string{
		filepath.Join(moduleRoot, "openapi", "dist", "openapi.json"),
		filepath.Join(moduleRoot, "openapi", "generated", "ts"),
		filepath.Join(moduleRoot, "openapi", "generated", "openapi-ts-tmp"),
	} {
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("清理生成目录失败 %s: %v", path, err)
		}
	}

	if matches, err := filepath.Glob(filepath.Join(moduleRoot, "..", "web", "openapi-ts-error-*.log")); err == nil {
		for _, path := range matches {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				t.Fatalf("清理 openapi-ts 错误日志失败 %s: %v", path, err)
			}
		}
	} else {
		t.Fatalf("扫描 openapi-ts 错误日志失败: %v", err)
	}

	ensureDistGitkeep(t, moduleRoot)
}

func assertBundleHasRequiredSurface(t *testing.T, bundlePath string) {
	t.Helper()

	raw, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatalf("读取 bundle 失败: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("解析 bundle 失败: %v", err)
	}

	assertBundlePaths(t, doc)
	assertBundleComponents(t, doc)
}

func assertBundlePaths(t *testing.T, doc map[string]any) {
	t.Helper()

	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 paths 节点")
	}

	for _, path := range []string{
		"/api/v1/public/sessions",
		"/api/v1/public/verification-tokens",
		"/api/v1/public/password-reset-requests",
		"/api/v1/public/password-resets",
		"/api/v1/user/me/sessions",
		"/api/v1/admin/users",
		"/api/v1/node/usage-reports",
	} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("bundle 中缺少关键 path: %s", path)
		}
	}
}

func assertBundleComponents(t *testing.T, doc map[string]any) {
	t.Helper()

	components, ok := doc["components"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 components 节点")
	}

	assertBundleComponentGroup(t, components, "parameters", []string{"page", "pageSize", "q", "sort", "order", "idempotencyKey", "nodeTimestamp", "nodeNonce", "nodeSignature"})
	assertBundleComponentGroup(t, components, "schemas", []string{"Problem", "ValidationProblem", "PaginationMeta", "Money"})
	assertBundleComponentGroup(t, components, "securitySchemes", []string{"sessionAuth", "nodeAuth"})
	assertNodeAuthStructure(t, components)
}

func assertNodeUsageReportOperation(t *testing.T, bundlePath string) {
	t.Helper()

	doc := readBundleDoc(t, bundlePath)
	paths := bundlePaths(t, doc)

	operation := bundleOperation(t, paths, "/api/v1/node/usage-reports", "post")
	assertOperationSecurity(t, operation, "nodeAuth")
	assertOperationParameters(t, operation, []string{
		"#/components/parameters/idempotencyKey",
		"#/components/parameters/nodeTimestamp",
		"#/components/parameters/nodeNonce",
		"#/components/parameters/nodeSignature",
	})
	assertUsageReportRequestBody(t, operation)
}

func readBundleDoc(t *testing.T, bundlePath string) map[string]any {
	t.Helper()

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

func bundlePaths(t *testing.T, doc map[string]any) map[string]any {
	t.Helper()

	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 paths 节点")
	}

	return paths
}

func bundleOperation(t *testing.T, paths map[string]any, pathName, method string) map[string]any {
	t.Helper()

	rawItem, ok := paths[pathName]
	if !ok {
		t.Fatalf("bundle 中缺少关键 path: %s", pathName)
	}

	item, ok := rawItem.(map[string]any)
	if !ok {
		t.Fatalf("路径 %s 的定义格式不正确", pathName)
	}

	rawOperation, ok := item[method]
	if !ok {
		t.Fatalf("路径 %s 缺少 %s 操作", pathName, method)
	}

	operation, ok := rawOperation.(map[string]any)
	if !ok {
		t.Fatalf("路径 %s 的 %s 操作格式不正确", pathName, method)
	}

	return operation
}

func assertOperationSecurity(t *testing.T, operation map[string]any, requiredScheme string) {
	t.Helper()

	rawSecurity, ok := operation["security"].([]any)
	if !ok {
		t.Fatalf("操作缺少 security 定义")
	}

	found := false
	for _, entry := range rawSecurity {
		schemes, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		if _, ok := schemes[requiredScheme]; ok {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("操作缺少 %s 安全方案", requiredScheme)
	}
}

func assertOperationParameters(t *testing.T, operation map[string]any, expectedRefs []string) {
	t.Helper()

	rawParams, ok := operation["parameters"].([]any)
	if !ok {
		t.Fatalf("操作缺少 parameters 定义")
	}
	if len(rawParams) != len(expectedRefs) {
		t.Fatalf("操作 parameters 数量不匹配，期望 %d，实际 %d", len(expectedRefs), len(rawParams))
	}

	got := make([]string, 0, len(rawParams))
	for _, rawParam := range rawParams {
		param, ok := rawParam.(map[string]any)
		if !ok {
			t.Fatalf("parameters 条目格式不正确")
		}
		ref, ok := param["$ref"].(string)
		if !ok {
			t.Fatalf("parameters 条目缺少 $ref")
		}
		got = append(got, normalizeComponentRef(ref))
	}

	for _, expected := range expectedRefs {
		if !containsString(got, normalizeComponentRef(expected)) {
			t.Fatalf("操作 parameters 缺少 %s，实际为 %v", expected, got)
		}
	}
}

func assertUsageReportRequestBody(t *testing.T, operation map[string]any) {
	t.Helper()

	requestBody, ok := operation["requestBody"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 缺少 requestBody")
	}
	content, ok := requestBody["content"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 缺少 requestBody.content")
	}
	jsonMedia, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 缺少 application/json 内容")
	}
	schema, ok := jsonMedia["schema"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 缺少 request schema")
	}
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 请求体缺少 properties")
	}
	traffic, ok := properties["traffic"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 请求体缺少 traffic")
	}
	trafficProps, ok := traffic["properties"].(map[string]any)
	if !ok {
		t.Fatalf("node usage report 请求体缺少 traffic.properties")
	}
	assertIntegerProperty(t, trafficProps, "uploadBytes")
	assertIntegerProperty(t, trafficProps, "downloadBytes")
}

func assertIntegerProperty(t *testing.T, properties map[string]any, name string) {
	t.Helper()

	rawProp, ok := properties[name]
	if !ok {
		t.Fatalf("缺少属性 %s", name)
	}
	prop, ok := rawProp.(map[string]any)
	if !ok {
		t.Fatalf("属性 %s 的定义格式不正确", name)
	}
	if typ, _ := prop["type"].(string); typ != "integer" {
		t.Fatalf("属性 %s 不是 integer 类型，实际为 %v", name, prop["type"])
	}
}

func assertNodeAuthStructure(t *testing.T, components map[string]any) {
	t.Helper()

	rawSchemes, ok := components["securitySchemes"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 components.securitySchemes")
	}

	rawNodeAuth, ok := rawSchemes["nodeAuth"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 nodeAuth 安全方案")
	}
	if typ, _ := rawNodeAuth["type"].(string); typ != "apiKey" {
		t.Fatalf("nodeAuth 类型不正确")
	}
	if location, _ := rawNodeAuth["in"].(string); location != "header" {
		t.Fatalf("nodeAuth 位置不正确")
	}
	if name, _ := rawNodeAuth["name"].(string); name != "X-Node-Key-Id" {
		t.Fatalf("nodeAuth key 名称不正确")
	}

	binding, ok := rawNodeAuth["x-node-auth-binding"].(map[string]any)
	if !ok {
		t.Fatalf("nodeAuth 缺少 x-node-auth-binding")
	}
	if scope, _ := binding["subjectScope"].(string); scope == "" {
		t.Fatalf("nodeAuth 缺少 subjectScope 约束")
	}
	if revocation, _ := binding["revocation"].(string); revocation == "" {
		t.Fatalf("nodeAuth 缺少 revocation 约束")
	}
	if skew, ok := rawNodeAuth["x-node-auth-clock-skew-seconds"].(float64); !ok || skew != 300 {
		t.Fatalf("nodeAuth 时间窗约束不正确")
	}
	if window, ok := rawNodeAuth["x-node-auth-nonce-reuse-window-seconds"].(float64); !ok || window != 300 {
		t.Fatalf("nodeAuth nonce 窗口约束不正确")
	}
	rawCoverage, ok := rawNodeAuth["x-node-auth-signature-coverage"].([]any)
	if !ok {
		t.Fatalf("nodeAuth 缺少签名覆盖范围")
	}
	coverage := make([]string, 0, len(rawCoverage))
	for _, item := range rawCoverage {
		text, ok := item.(string)
		if !ok {
			t.Fatalf("nodeAuth 签名覆盖范围格式不正确")
		}
		coverage = append(coverage, text)
	}
	for _, expected := range []string{"method", "normalizedPath", "timestamp", "nonce", "bodyDigest"} {
		if !containsString(coverage, expected) {
			t.Fatalf("nodeAuth 签名覆盖范围缺少 %s", expected)
		}
	}
}

func assertBundleComponentGroup(t *testing.T, components map[string]any, group string, required []string) {
	t.Helper()

	rawGroup, ok := components[group].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 components.%s", group)
	}

	for _, name := range required {
		if _, ok := rawGroup[name]; !ok {
			t.Fatalf("bundle 中缺少 components.%s.%s", group, name)
		}
	}
}

func assertPathFilesUseRootComponents(t *testing.T, moduleRoot string) {
	t.Helper()

	pathDir := filepath.Join(moduleRoot, "openapi", "paths")
	err := filepath.WalkDir(pathDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		content := string(raw)
		if strings.Contains(content, "../../components/") {
			t.Fatalf("路径文件仍直接引用底层 components: %s", path)
		}
		if strings.Contains(content, "#/components/") && !strings.Contains(content, "../../openapi.yaml#/components/") {
			t.Fatalf("路径文件没有通过根 openapi.yaml 引用公共组件: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("扫描 path 文件失败: %v", err)
	}

	rootSpec, err := os.ReadFile(filepath.Join(moduleRoot, "openapi", "openapi.yaml"))
	if err != nil {
		t.Fatalf("读取根 OpenAPI 失败: %v", err)
	}
	if strings.Contains(string(rootSpec), "\n  headers:\n") {
		t.Fatalf("根 OpenAPI 仍然暴露 headers 组件分组")
	}
}

func assertRootSpecUsesComponentRefs(t *testing.T, moduleRoot string) {
	t.Helper()

	rootSpec, err := os.ReadFile(filepath.Join(moduleRoot, "openapi", "openapi.yaml"))
	if err != nil {
		t.Fatalf("读取根 OpenAPI 失败: %v", err)
	}

	content := string(rootSpec)
	for _, required := range []string{
		"$ref: ./components/parameters/page.yaml",
		"$ref: ./components/parameters/idempotency_key.yaml",
		"$ref: ./components/schemas/common/problem.yaml",
		"$ref: ./components/schemas/common/validation_problem.yaml",
		"$ref: ./components/security/security.yaml#/nodeAuth",
	} {
		if !strings.Contains(content, required) {
			t.Fatalf("根 OpenAPI 缺少组件引用: %s", required)
		}
	}
}

func assertGeneratedTypesHaveNoDuplicatePublicAliases(t *testing.T, generatedPath string) {
	t.Helper()

	raw, err := os.ReadFile(filepath.Join(generatedPath, "types.gen.ts"))
	if err != nil {
		t.Fatalf("读取生成的类型文件失败: %v", err)
	}

	content := string(raw)
	for _, forbidden := range []string{"Problem2", "ValidationProblem2", "PaginationMeta2", "Money2"} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("生成产物仍包含公共组件重复别名: %s", forbidden)
		}
	}
}

func assertNoOpenApiTsErrorLogs(t *testing.T, moduleRoot string) {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(moduleRoot, "..", "web", "openapi-ts-error-*.log"))
	if err != nil {
		t.Fatalf("扫描 openapi-ts 错误日志失败: %v", err)
	}
	if len(matches) > 0 {
		t.Fatalf("仍存在 openapi-ts 错误日志: %v", matches)
	}
}

func ensureDistGitkeep(t *testing.T, moduleRoot string) {
	t.Helper()

	path := filepath.Join(moduleRoot, "openapi", "dist", ".gitkeep")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("恢复 dist 目录失败: %v", err)
	}
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("恢复 dist/.gitkeep 失败: %v", err)
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func normalizeComponentRef(ref string) string {
	if idx := strings.LastIndex(ref, "/"); idx >= 0 {
		ref = ref[idx+1:]
	}
	return normalizeComponentName(ref)
}

func normalizeComponentName(name string) string {
	replacer := strings.NewReplacer("_", "", "-", "")
	return strings.ToLower(replacer.Replace(name))
}
