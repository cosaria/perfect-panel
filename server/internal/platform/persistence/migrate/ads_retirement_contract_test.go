package migrate

import (
	"os"
	"strings"
	"testing"
)

func TestAdsRetirementMigrationExists(t *testing.T) {
	up, err := os.ReadFile("database/02126_remove_ads.up.sql")
	if err != nil {
		t.Fatalf("expected ads retirement up migration: %v", err)
	}
	down, err := os.ReadFile("database/02126_remove_ads.down.sql")
	if err != nil {
		t.Fatalf("expected ads retirement down migration: %v", err)
	}

	upSQL := string(up)
	downSQL := string(down)

	if !strings.Contains(upSQL, "DROP TABLE IF EXISTS `ads`") {
		t.Fatal("expected ads retirement migration to drop ads table")
	}
	if !strings.Contains(upSQL, "show_ads") {
		t.Fatal("expected ads retirement migration to clean device show_ads config")
	}
	if !strings.Contains(downSQL, "CREATE TABLE IF NOT EXISTS `ads`") {
		t.Fatal("expected ads retirement down migration to recreate ads table")
	}
}
