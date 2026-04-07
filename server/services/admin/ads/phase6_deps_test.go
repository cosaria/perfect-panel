package ads

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase6AdminAdsUsesPackageLocalDeps(t *testing.T) {
	_, err := os.Stat("deps.go")
	require.NoError(t, err)

	err = filepath.Walk(".", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", path)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", path)
		}
		return nil
	})
	require.NoError(t, err)
}
