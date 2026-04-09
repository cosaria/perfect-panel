package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase3ServicesNoLongerKeepHandlerLogicPairs(t *testing.T) {
	servicesRoot := filepath.Join("..", "internal", "domains")
	legacyPairs := 0

	err := filepath.Walk(servicesRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(path, "Handler.go") {
			return nil
		}

		logicPath := strings.TrimSuffix(path, "Handler.go") + "Logic.go"
		if _, err := os.Stat(logicPath); err == nil {
			legacyPairs++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk domains tree: %v", err)
	}

	if legacyPairs != 0 {
		t.Fatalf("expected no split handler/logic pairs in domains, found %d", legacyPairs)
	}
}

func TestPhase4WorkerNoLongerKeepsLegacySubpackages(t *testing.T) {
	for _, legacyDir := range []string{
		filepath.Join("..", "worker", "handler"),
		filepath.Join("..", "worker", "logic"),
		filepath.Join("..", "worker", "types"),
	} {
		if _, err := os.Stat(legacyDir); err == nil {
			t.Fatalf("expected legacy worker directory %s to be removed", legacyDir)
		} else if !os.IsNotExist(err) {
			t.Fatalf("stat %s: %v", legacyDir, err)
		}
	}
}
