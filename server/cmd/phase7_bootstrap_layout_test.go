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

func TestPhase7LegacyBootstrapDirectoriesRemoved(t *testing.T) {
	legacy := []string{
		filepath.Join("..", "svc"),
		filepath.Join("..", "initialize"),
		filepath.Join("..", "runtime"),
	}

	for _, target := range legacy {
		if _, err := os.Stat(target); err == nil {
			t.Fatalf("expected legacy bootstrap directory %s to be removed", target)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", target, err)
		}
	}
}
