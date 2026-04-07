package common

import (
	"testing"

	"github.com/perfect-panel/server/config"
)

func TestPhase6DepsExposeLiveConfigPointer(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{}
	cfg.Site.SiteName = "before"

	deps := Deps{Config: cfg}
	cfg.Site.SiteName = "after"

	if deps.currentConfig().Site.SiteName != "after" {
		t.Fatalf("expected live config view, got %q", deps.currentConfig().Site.SiteName)
	}
}
