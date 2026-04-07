package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase6OpenAPICommandWritesSpecs(t *testing.T) {
	outputDir := t.TempDir()
	if err := openapiCmd.Flags().Set("output", outputDir); err != nil {
		t.Fatalf("set output flag: %v", err)
	}
	t.Cleanup(func() {
		_ = openapiCmd.Flags().Set("output", "docs/openapi")
	})

	if err := openapiCmd.RunE(openapiCmd, nil); err != nil {
		t.Fatalf("run openapi command: %v", err)
	}

	for _, name := range []string{"admin.json", "common.json", "user.json"} {
		path := filepath.Join(outputDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected exported spec %s: %v", path, err)
		}
	}
}
