package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7HTTPDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "platform", "http", "routes.go"),
		filepath.Join("..", "internal", "platform", "http", "middleware", "authMiddleware.go"),
		filepath.Join("..", "internal", "platform", "http", "response", "response.go"),
		filepath.Join("..", "internal", "platform", "http", "types", "types.go"),
		filepath.Join("..", "cmd", "openapi", "main.go"),
	}

	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected HTTP target %s to exist: %v", target, err)
		}
	}
}
