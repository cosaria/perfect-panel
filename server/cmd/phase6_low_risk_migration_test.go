package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase6LowRiskBatchNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "admin", "announcement"),
		filepath.Join("..", "services", "admin", "coupon"),
		filepath.Join("..", "services", "admin", "document"),
		filepath.Join("..", "services", "common", "getAds.go"),
		filepath.Join("..", "services", "common", "getClient.go"),
		filepath.Join("..", "services", "user", "announcement"),
		filepath.Join("..", "services", "user", "document"),
		filepath.Join("..", "services", "user", "payment", "getAvailablePaymentMethods.go"),
		filepath.Join("..", "services", "user", "portal", "getAvailablePaymentMethods.go"),
	}

	for _, target := range targets {
		assertNoServiceContextDependency(t, target)
	}
}

func assertNoServiceContextDependency(t *testing.T, target string) {
	t.Helper()

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat %s: %v", target, err)
	}

	if !info.IsDir() {
		assertFileHasNoServiceContextDependency(t, target)
		return
	}

	err = filepath.Walk(target, func(path string, walkInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if walkInfo.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		assertFileHasNoServiceContextDependency(t, path)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", target, err)
	}
}

func assertFileHasNoServiceContextDependency(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	source := string(content)
	if strings.Contains(source, "*svc.ServiceContext") {
		t.Fatalf("%s still depends on *svc.ServiceContext", path)
	}
	if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
		t.Fatalf("%s still imports server/svc", path)
	}
}
