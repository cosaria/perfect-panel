package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/internal/platform/http"
	"github.com/stretchr/testify/require"
)

func TestPhase8DatabaseNormalizationContracts(t *testing.T) {
	contentTargets := []string{
		filepath.Join("..", "internal", "platform", "persistence", "content", "models.go"),
		filepath.Join("..", "internal", "platform", "persistence", "content", "repository.go"),
	}
	for _, target := range contentTargets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected content normalization target %s to exist: %v", target, err)
		}
	}

	disallowMigrateImports := []string{
		filepath.Join("..", "internal", "bootstrap", "configinit"),
		filepath.Join("..", "internal", "bootstrap", "app"),
		filepath.Join("..", "cmd"),
	}
	for _, root := range disallowMigrateImports {
		assertRootHasNoMigrateDependency(t, root)
	}

	for _, target := range []string{
		filepath.Join("..", "internal", "domains", "user", "ticket"),
		filepath.Join("..", "internal", "domains", "user", "document"),
		filepath.Join("..", "internal", "domains", "user", "announcement"),
		filepath.Join("..", "internal", "domains", "node", "serverPushUserTraffic.go"),
	} {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected contract target %s to exist: %v", target, err)
		}
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	apis := handler.RegisterHandlersForSpec(router)
	userSpec, err := apis.UserOpenAPI()
	require.NoError(t, err)

	for _, spec := range []interface{}{apis.Admin.OpenAPI(), apis.Common.OpenAPI(), userSpec} {
		data, err := json.Marshal(spec)
		require.NoError(t, err)
		require.Contains(t, string(data), "\"openapi\"")
	}
}

func assertRootHasNoMigrateDependency(t *testing.T, root string) {
	t.Helper()

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		source := string(content)
		if strings.Contains(source, "internal/platform/persistence/migrate") || strings.Contains(source, "migrate.") {
			t.Fatalf("%s should not depend on legacy migrate runner", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
}
