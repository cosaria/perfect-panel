package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase6MixedDependencyCheckpointNoLongerImportsServiceContext(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "common", "sendEmailCode.go"),
		filepath.Join("..", "services", "common", "sendSmsCode.go"),
		filepath.Join("..", "services", "common", "checkverificationcodehandler.go"),
		filepath.Join("..", "services", "common", "checkverificationcodelogic.go"),
	}

	for _, target := range targets {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}

		source := string(content)
		if strings.Contains(source, "*svc.ServiceContext") {
			t.Fatalf("%s still depends on *svc.ServiceContext", target)
		}
		if strings.Contains(source, "\"github.com/perfect-panel/server/svc\"") {
			t.Fatalf("%s still imports server/svc", target)
		}
	}
}
