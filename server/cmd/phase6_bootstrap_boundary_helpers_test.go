package cmd_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	phase6ServiceContextUsagePattern = regexp.MustCompile(`\b(?:svc|appbootstrap)\.ServiceContext\b`)
	phase6BootstrapImportPattern     = regexp.MustCompile(`"github\.com/perfect-panel/server/(?:svc|initialize|runtime|internal/bootstrap/(?:app|configinit|runtime))"`)
)

func assertTargetsHaveNoBootstrapBoundaryDependency(t *testing.T, targets []string) {
	t.Helper()

	for _, target := range targets {
		assertNoBootstrapBoundaryDependency(t, target)
	}
}

func assertNoBootstrapBoundaryDependency(t *testing.T, target string) {
	t.Helper()

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat %s: %v", target, err)
	}

	if !info.IsDir() {
		assertFileHasNoBootstrapBoundaryDependency(t, target)
		return
	}

	err = filepath.Walk(target, func(path string, walkInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if walkInfo.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		assertFileHasNoBootstrapBoundaryDependency(t, path)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", target, err)
	}
}

func assertFileHasNoBootstrapBoundaryDependency(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	source := string(content)
	if match := phase6ServiceContextUsagePattern.FindString(source); match != "" {
		t.Fatalf("%s still depends on bootstrap ServiceContext alias %q", path, match)
	}
	if match := phase6BootstrapImportPattern.FindString(source); match != "" {
		t.Fatalf("%s still imports bootstrap composition root %s", path, match)
	}
}

func globPhase6Targets(t *testing.T, pattern string) []string {
	t.Helper()

	targets, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob %s: %v", pattern, err)
	}
	return targets
}

func countPatternInGoFiles(t *testing.T, root string, pattern *regexp.Regexp) int {
	t.Helper()

	count := 0
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

		count += len(pattern.FindAll(content, -1))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return count
}
