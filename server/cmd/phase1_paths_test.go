package cmd_test

import (
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	servermigrate "github.com/perfect-panel/server/models/migrate"
	serversyncx "github.com/perfect-panel/server/modules/cache/syncx"
	serverrules "github.com/perfect-panel/server/modules/util/rules"
	serversvc "github.com/perfect-panel/server/svc"
	servertypes "github.com/perfect-panel/server/types"
)

func TestPhase1TopLevelPathsExist(t *testing.T) {
	var cfg serverconfig.Config
	var ctx serversvc.ServiceContext
	var ads servertypes.Ads
	limit := serversyncx.NewLimit(1)
	rule := serverrules.NewRule("DOMAIN,example.com,DIRECT", "fallback")

	if servermigrate.NoChange == nil {
		t.Fatal("expected migrate package to expose NoChange")
	}

	if cfg.Port != 0 {
		t.Fatal("expected zero-value config in compile smoke test")
	}

	if serverconfig.Version == "" || serverconfig.BuildTime == "" || serverconfig.Repository == "" || serverconfig.ServiceName == "" {
		t.Fatal("expected config package to expose build metadata")
	}

	if ctx.Redis != nil {
		t.Fatal("expected zero-value service context in compile smoke test")
	}

	if ads.Id != 0 {
		t.Fatal("expected zero-value type in compile smoke test")
	}

	if !limit.TryBorrow() {
		t.Fatal("expected syncx limit from modules/cache to be usable")
	}

	if rule == nil {
		t.Fatal("expected rules package from modules/util to parse a basic rule")
	}
}
