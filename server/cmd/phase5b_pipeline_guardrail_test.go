package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhase5BOpenAPIGovernanceWorkflowExists(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "openapi-governance.yml")
	_, err := os.Stat(path)
	require.NoError(t, err)
}
