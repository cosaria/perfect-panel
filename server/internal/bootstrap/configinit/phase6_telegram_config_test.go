package configinit

import (
	"testing"

	"github.com/perfect-panel/server/models/auth"
)

func TestPhase6TelegramRuntimeConfigUsesAuthConfigValues(t *testing.T) {
	t.Parallel()

	cfg := runtimeTelegramConfig(&auth.TelegramAuthConfig{
		BotToken:      "telegram-bot-token",
		EnableNotify:  true,
		WebHookDomain: "https://example.com",
	})

	if cfg.BotToken != "telegram-bot-token" {
		t.Fatalf("expected bot token to come from auth config, got %q", cfg.BotToken)
	}
	if !cfg.EnableNotify {
		t.Fatal("expected enable notify to come from auth config")
	}
	if cfg.WebHookDomain != "https://example.com" {
		t.Fatalf("expected webhook domain to come from auth config, got %q", cfg.WebHookDomain)
	}
}
