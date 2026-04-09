package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6MixedDependencyCheckpointNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	assertTargetsHaveNoBootstrapBoundaryDependency(t, []string{
		filepath.Join("..", "services", "common", "sendEmailCode.go"),
		filepath.Join("..", "services", "common", "sendSmsCode.go"),
		filepath.Join("..", "services", "common", "checkverificationcodehandler.go"),
		filepath.Join("..", "services", "common", "checkverificationcodelogic.go"),
	})
}
