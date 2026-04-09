package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7DomainDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "domains", "admin"),
		filepath.Join("..", "internal", "domains", "auth"),
		filepath.Join("..", "internal", "domains", "common"),
		filepath.Join("..", "internal", "domains", "common", "report"),
		filepath.Join("..", "internal", "domains", "node"),
		filepath.Join("..", "internal", "domains", "subscribe"),
		filepath.Join("..", "internal", "domains", "telegram"),
		filepath.Join("..", "internal", "domains", "user"),
	}

	for _, target := range targets {
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("expected domain target %s to exist: %v", target, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %s to be a directory", target)
		}
	}
}

func TestPhase7LegacyServicesDirectoriesRemoved(t *testing.T) {
	legacy := []string{
		filepath.Join("..", "services"),
		filepath.Join("..", "services", "admin"),
		filepath.Join("..", "services", "auth"),
		filepath.Join("..", "services", "common"),
		filepath.Join("..", "services", "node"),
		filepath.Join("..", "services", "notify"),
		filepath.Join("..", "services", "report"),
		filepath.Join("..", "services", "subscribe"),
		filepath.Join("..", "services", "telegram"),
		filepath.Join("..", "services", "user"),
	}

	for _, target := range legacy {
		if _, err := os.Stat(target); err == nil {
			t.Fatalf("expected legacy services directory %s to be removed", target)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", target, err)
		}
	}
}
