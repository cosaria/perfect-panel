package cmd_test

import (
	"os"
	"os/exec"
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

func TestPhase7RootEntryBuilds(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "ppanel-test-bin")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join("..")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected server root entry to build, got err=%v\n%s", err, string(output))
	}
	if _, err := os.Stat(binPath); err != nil {
		t.Fatalf("expected built binary at %s: %v", binPath, err)
	}
}

func TestPhase7DockerfilesBuildFromModuleRoot(t *testing.T) {
	dockerfiles := []struct {
		name string
		path string
	}{
		{name: "repo-root Dockerfile", path: filepath.Join("..", "..", "Dockerfile")},
		{name: "server Dockerfile", path: filepath.Join("..", "Dockerfile")},
	}

	for _, item := range dockerfiles {
		data, err := os.ReadFile(item.path)
		if err != nil {
			t.Fatalf("read %s: %v", item.name, err)
		}
		source := string(data)
		if strings.Contains(source, "ppanel.go") {
			t.Fatalf("expected %s to avoid legacy ppanel.go build entry, got:\n%s", item.name, source)
		}
		if !strings.Contains(source, "go build") {
			t.Fatalf("expected %s to build Go entrypoint, got:\n%s", item.name, source)
		}
		if !strings.Contains(source, "-o /app/ppanel .") {
			t.Fatalf("expected %s to build from module root with '.', got:\n%s", item.name, source)
		}
	}
}
