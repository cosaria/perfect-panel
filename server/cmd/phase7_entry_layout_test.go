package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase7RootMainDelegatesToCmdPpanel(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	source := string(data)
	if !strings.Contains(source, "\"github.com/perfect-panel/server/cmd/ppanel\"") {
		t.Fatalf("expected root main to import cmd/ppanel, got:\n%s", source)
	}
	if !strings.Contains(source, "ppanel.Execute()") {
		t.Fatalf("expected root main to delegate to ppanel.Execute, got:\n%s", source)
	}
}

func TestPhase7PpanelCommandPackageExists(t *testing.T) {
	if _, err := os.Stat(filepath.Join("ppanel", "root.go")); err != nil {
		t.Fatalf("expected cmd/ppanel/root.go to exist: %v", err)
	}
}
