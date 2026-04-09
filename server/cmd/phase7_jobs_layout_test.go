package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7JobsDirectoryExists(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "jobs", "consumer_service.go"),
		filepath.Join("..", "internal", "jobs", "scheduler_service.go"),
		filepath.Join("..", "internal", "jobs", "registry", "routes.go"),
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected jobs target %s to exist: %v", target, err)
		}
	}
}

func TestPhase7LegacyJobsDirectoriesRemoved(t *testing.T) {
	legacy := []string{
		filepath.Join("..", "worker"),
	}

	for _, target := range legacy {
		if _, err := os.Stat(target); err == nil {
			t.Fatalf("expected legacy jobs directory %s to be removed", target)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", target, err)
		}
	}
}
