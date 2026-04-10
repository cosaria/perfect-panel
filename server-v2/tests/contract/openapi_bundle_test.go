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
	assertPathFilesUseRootComponents(t, moduleRoot)
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
