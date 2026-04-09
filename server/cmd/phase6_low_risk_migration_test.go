package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6LowRiskBatchNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "internal", "domains", "admin", "announcement"),
		filepath.Join("..", "internal", "domains", "admin", "coupon"),
		filepath.Join("..", "internal", "domains", "admin", "document"),
		filepath.Join("..", "internal", "domains", "common", "getClient.go"),
		filepath.Join("..", "internal", "domains", "user", "announcement"),
		filepath.Join("..", "internal", "domains", "user", "document"),
		filepath.Join("..", "internal", "domains", "user", "payment", "getAvailablePaymentMethods.go"),
		filepath.Join("..", "internal", "domains", "user", "portal", "getAvailablePaymentMethods.go"),
	}

	assertTargetsHaveNoBootstrapBoundaryDependency(t, targets)
}
