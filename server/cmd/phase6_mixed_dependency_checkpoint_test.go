package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6MixedDependencyCheckpointNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "internal", "domains", "common", "sendEmailCode.go"),
		filepath.Join("..", "internal", "domains", "common", "sendSmsCode.go"),
		filepath.Join("..", "internal", "domains", "common", "checkverificationcodehandler.go"),
		filepath.Join("..", "internal", "domains", "common", "checkverificationcodelogic.go"),
	})
}
