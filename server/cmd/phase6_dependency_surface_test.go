package cmd_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var serviceContextUsagePattern = regexp.MustCompile(`\*svc\.ServiceContext`)

func TestPhase6ServiceContextExplainsCompositionRootRole(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "svc", "serviceContext.go"))
	require.NoError(t, err)

	require.Contains(t, string(content), "ServiceContext is a temporary composition-root shell")
}

func TestPhase6DependencySurfaceBaseline(t *testing.T) {
	t.Parallel()

	baselineMaximum := map[string]int{
		filepath.Join("..", "services"):              738,
		filepath.Join("..", "worker"):                31,
		filepath.Join("..", "routers", "middleware"): 8,
		filepath.Join("..", "initialize"):            14,
	}

	total := 0
	for root, maxAllowed := range baselineMaximum {
		got := countServiceContextUsages(t, root)
		if got > maxAllowed {
			t.Fatalf("expected %s to contain at most %d direct ServiceContext usages, got %d", root, maxAllowed, got)
		}
		total += got
	}

	require.LessOrEqual(t, total, 791)
}

func countServiceContextUsages(t *testing.T, root string) int {
	t.Helper()

	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		count += len(serviceContextUsagePattern.FindAll(content, -1))
		return nil
	})
	require.NoError(t, err)
	return count
}
