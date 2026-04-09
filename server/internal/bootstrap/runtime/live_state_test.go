package runtime

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
)

func TestLiveStateTracksMutableRuntimeValues(t *testing.T) {
	t.Parallel()

	state := NewLiveState()

	state.SetExchangeRate(6.66)
	if got := state.ExchangeRate(); got != 6.66 {
		t.Fatalf("expected exchange rate 6.66, got %v", got)
	}

	called := 0
	state.SetRestart(func() error {
		called++
		return nil
	})
	restart := state.Restart()
	if restart == nil {
		t.Fatal("expected restart capability to be available")
	}
	if err := restart(); err != nil {
		t.Fatalf("restart failed: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected restart to be called once, got %d", called)
	}

	bot := &tgbotapi.BotAPI{}
	state.SetTelegramBot(bot)
	if got := state.TelegramBot(); got != bot {
		t.Fatal("expected telegram bot pointer to round-trip through live state")
	}

	manager := &node.Manager{}
	state.SetNodeMultiplierManager(manager)
	if got := state.NodeMultiplierManager(); got != manager {
		t.Fatal("expected node multiplier manager to round-trip through live state")
	}
}

func TestLiveStatePreventsStaleExchangeRateWritesAcrossCurrencyChanges(t *testing.T) {
	t.Parallel()

	state := NewLiveState()

	version := state.PrepareExchangeRate("USD", "CNY")
	if !state.StoreExchangeRate(version, "USD", "CNY", 7.11) {
		t.Fatal("expected initial exchange rate write to succeed")
	}

	quote := state.ExchangeRateQuote()
	if quote.From != "USD" || quote.To != "CNY" || quote.Rate != 7.11 {
		t.Fatalf("unexpected initial quote: %+v", quote)
	}

	newVersion := state.PrepareExchangeRate("EUR", "CNY")
	if newVersion == version {
		t.Fatal("expected currency change to bump exchange rate version")
	}
	if state.StoreExchangeRate(version, "USD", "CNY", 7.22) {
		t.Fatal("expected stale exchange rate write to be rejected")
	}

	quote = state.ExchangeRateQuote()
	if quote.From != "EUR" || quote.To != "CNY" {
		t.Fatalf("expected quote pair to track latest currency, got %+v", quote)
	}
	if quote.Rate != 0 {
		t.Fatalf("expected currency switch to invalidate cached rate, got %+v", quote)
	}
}

func TestLiveStateSwapsTelegramPollers(t *testing.T) {
	t.Parallel()

	state := NewLiveState()
	stopped := 0

	if previous := state.SwapTelegramPoller(func() { stopped++ }); previous != nil {
		t.Fatal("expected first poller registration to have no previous stop func")
	}

	previous := state.SwapTelegramPoller(func() { stopped += 10 })
	if previous == nil {
		t.Fatal("expected second poller registration to return previous stop func")
	}
	previous()
	if stopped != 1 {
		t.Fatalf("expected previous poller stop func to be callable, got %d", stopped)
	}
}
