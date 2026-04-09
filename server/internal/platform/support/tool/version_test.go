package tool

import (
	"testing"

	"github.com/perfect-panel/server/config"
)

func TestExtractVersionNumber(t *testing.T) {
	versionNumber := ExtractVersionNumber(config.Version)
	t.Log(versionNumber)
}
