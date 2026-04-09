package handler

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/models/node"
)

func TestPhase6DynamicRouteDepsFollowLiveState(t *testing.T) {
	t.Parallel()

	runtimeDeps := &appruntime.Deps{
		Config: &config.Config{},
		Live:   appruntime.NewLiveState(),
	}
	initDeps := initializeDepsFromRuntimeDeps(runtimeDeps)

	publicPortalDeps := newPublicPortalDeps(runtimeDeps)
	publicUserDeps := newPublicUserDeps(runtimeDeps)
	adminSystemDeps := newAdminSystemDeps(runtimeDeps, initDeps)
	adminToolDeps := newAdminToolDeps(runtimeDeps)
	telegramDeps := newTelegramServiceDeps(runtimeDeps)

	bot := &tgbotapi.BotAPI{}
	manager := &node.Manager{}
	restarted := false

	runtimeDeps.Live.SetExchangeRate(1.23)
	runtimeDeps.Live.SetTelegramBot(bot)
	runtimeDeps.Live.SetNodeMultiplierManager(manager)
	runtimeDeps.Live.SetRestart(func() error {
		restarted = true
		return nil
	})

	if got := publicPortalDeps.CurrentExchangeRate(); got != 1.23 {
		t.Fatalf("expected portal deps to read live exchange rate, got %v", got)
	}
	version := runtimeDeps.Live.PrepareExchangeRate("USD", "CNY")
	if !publicPortalDeps.StoreExchangeRateCache(version, "USD", "CNY", 6.54) {
		t.Fatal("expected portal deps to write versioned exchange rate cache")
	}
	quote := publicPortalDeps.CurrentExchangeRateSnapshot()
	if quote.From != "USD" || quote.To != "CNY" || quote.Rate != 6.54 {
		t.Fatalf("expected portal deps to expose versioned exchange rate snapshot, got %+v", quote)
	}
	newVersion := publicPortalDeps.PrepareExchangeRateCache("EUR", "CNY")
	if publicPortalDeps.StoreExchangeRateCache(version, "USD", "CNY", 7.77) {
		t.Fatal("expected stale portal cache write to be rejected after currency change")
	}
	if !publicPortalDeps.StoreExchangeRateCache(newVersion, "EUR", "CNY", 8.88) {
		t.Fatal("expected fresh portal cache write to succeed")
	}
	if got := publicUserDeps.CurrentTelegramBot(); got != bot {
		t.Fatal("expected public user deps to read live telegram bot")
	}
	if got := telegramDeps.CurrentTelegramBot(); got != bot {
		t.Fatal("expected telegram route deps to read live telegram bot")
	}
	if got := adminSystemDeps.CurrentNodeMultiplierManager(); got != manager {
		t.Fatal("expected admin system deps to read live node multiplier manager")
	}
	if err := adminToolDeps.Restart(); err != nil {
		t.Fatalf("restart through admin tool deps failed: %v", err)
	}
	if !restarted {
		t.Fatal("expected restart closure to use live runtime capability")
	}
}
