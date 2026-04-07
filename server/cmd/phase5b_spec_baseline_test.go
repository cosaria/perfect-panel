package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase5BGovernanceDocumentExistsAndListsExclusions(t *testing.T) {
	path := filepath.Join("..", "..", "docs", "api-governance.md")

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	body := string(data)
	require.Contains(t, body, "OpenAPI Governance")
	require.Contains(t, body, "Excluded Surfaces")
	require.Contains(t, body, "node polling")
	require.Contains(t, body, "webhook")
	require.Contains(t, body, "init/bootstrap")
}

func TestPhase5BExportedSpecsDeclareGovernedTopLevelTags(t *testing.T) {
	specs := exportPhase5BSpecs(t)

	expectedTags := map[string][]string{
		"admin":  {"ads", "announcement", "application", "auth-method"},
		"common": {"common"},
		"user":   {"auth", "oauth", "order", "portal", "user"},
	}

	for name, expected := range expectedTags {
		t.Run(name, func(t *testing.T) {
			tags, ok := specs[name]["tags"].([]interface{})
			require.True(t, ok, "expected top-level tags in %s spec", name)
			require.NotEmpty(t, tags, "expected non-empty top-level tags in %s spec", name)

			describedTags := map[string]string{}
			for _, raw := range tags {
				tag, ok := raw.(map[string]interface{})
				require.True(t, ok)
				tagName, _ := tag["name"].(string)
				tagDescription, _ := tag["description"].(string)
				describedTags[tagName] = tagDescription
			}

			for _, tagName := range expected {
				require.Contains(t, describedTags, tagName)
				require.NotEmpty(t, strings.TrimSpace(describedTags[tagName]))
			}
		})
	}
}
