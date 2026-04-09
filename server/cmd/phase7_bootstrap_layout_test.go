package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7BootstrapDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "bootstrap", "app", "serviceContext.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "config.go"),
		filepath.Join("..", "internal", "bootstrap", "runtime", "live_state.go"),
	}

	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected bootstrap target %s to exist: %v", target, err)
		}
	}
}
