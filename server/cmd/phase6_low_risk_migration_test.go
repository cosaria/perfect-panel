package cmd_test

import (
	"path/filepath"
	"testing"
)

func TestPhase6LowRiskBatchNoLongerImportsBootstrapCompositionRoot(t *testing.T) {
	t.Parallel()

	targets := []string{
		filepath.Join("..", "services", "admin", "announcement"),
		filepath.Join("..", "services", "admin", "coupon"),
		filepath.Join("..", "services", "admin", "document"),
		filepath.Join("..", "services", "common", "getClient.go"),
		filepath.Join("..", "services", "user", "announcement"),
		filepath.Join("..", "services", "user", "document"),
		filepath.Join("..", "services", "user", "payment", "getAvailablePaymentMethods.go"),
		filepath.Join("..", "services", "user", "portal", "getAvailablePaymentMethods.go"),
	}

	assertTargetsHaveNoBootstrapBoundaryDependency(t, targets)
}
