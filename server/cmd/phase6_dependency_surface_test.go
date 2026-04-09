package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase6ServiceContextExplainsCompositionRootRole(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "internal", "bootstrap", "app", "serviceContext.go"))
	require.NoError(t, err)

	require.Contains(t, string(content), "ServiceContext is a temporary composition-root shell")
}

func TestPhase6DependencySurfaceBaseline(t *testing.T) {
	t.Parallel()

	serviceContextRoots := map[string]int{
		filepath.Join("..", "internal", "domains"):                        0,
		filepath.Join("..", "internal", "jobs"):                           0,
		filepath.Join("..", "internal", "platform", "http", "middleware"): 0,
		filepath.Join("..", "internal", "bootstrap", "app"):               0,
		filepath.Join("..", "internal", "bootstrap", "configinit"):        0,
		filepath.Join("..", "internal", "bootstrap", "runtime"):           0,
	}

	totalServiceContextUsages := 0
	for root, maxAllowed := range serviceContextRoots {
		got := countPatternInGoFiles(t, root, phase6ServiceContextUsagePattern)
		if got > maxAllowed {
			t.Fatalf("expected %s to contain at most %d bootstrap ServiceContext references, got %d", root, maxAllowed, got)
		}
		totalServiceContextUsages += got
	}
	require.Zero(t, totalServiceContextUsages)

	bootstrapImportRoots := map[string]int{
		filepath.Join("..", "internal", "domains"):                        0,
		filepath.Join("..", "internal", "jobs"):                           0,
		filepath.Join("..", "internal", "platform", "http", "middleware"): 8,
		filepath.Join("..", "internal", "bootstrap", "app"):               0,
		filepath.Join("..", "internal", "bootstrap", "configinit"):        0,
		filepath.Join("..", "internal", "bootstrap", "runtime"):           0,
	}

	totalBootstrapImports := 0
	for root, maxAllowed := range bootstrapImportRoots {
		got := countPatternInGoFiles(t, root, phase6BootstrapImportPattern)
		if got > maxAllowed {
			t.Fatalf("expected %s to contain at most %d bootstrap composition-root imports, got %d", root, maxAllowed, got)
		}
		totalBootstrapImports += got
	}
	require.LessOrEqual(t, totalBootstrapImports, 8)
}
