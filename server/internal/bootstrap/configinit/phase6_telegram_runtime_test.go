package configinit

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/models/auth"
	"gorm.io/gorm"
)

func TestSyncTelegramRuntimeUpdatesConfigAndClearsOldBot(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{}
	var currentBot *tgbotapi.BotAPI

	deps := Deps{
		Config: cfg,
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			currentBot = bot
		},
	}

	bot := &tgbotapi.BotAPI{}
	syncTelegramRuntime(deps, config.Telegram{
		BotToken:      "new-token",
		EnableNotify:  true,
		WebHookDomain: "https://example.com",
	}, bot)

	if cfg.Telegram.BotToken != "new-token" {
		t.Fatalf("expected runtime config token to update, got %q", cfg.Telegram.BotToken)
	}
	if currentBot != bot {
		t.Fatal("expected runtime bot to update on successful reload")
	}

	syncTelegramRuntime(deps, config.Telegram{}, nil)
	if cfg.Telegram.BotToken != "" {
		t.Fatalf("expected runtime config token to clear when telegram is disabled, got %q", cfg.Telegram.BotToken)
	}
	if currentBot != nil {
		t.Fatal("expected runtime bot to clear when telegram is disabled")
	}
}

func TestReplaceTelegramPollerStopsPreviousPoller(t *testing.T) {
	t.Parallel()

	stopped := 0
	var current func()
	deps := Deps{
		SwapTelegramPoller: func(next func()) func() {
			previous := current
			current = next
			return previous
		},
	}

	replaceTelegramPoller(deps, func() { stopped++ })
	replaceTelegramPoller(deps, func() { stopped += 10 })
	if stopped != 1 {
		t.Fatalf("expected previous poller to stop before replacement, got %d", stopped)
	}

	replaceTelegramPoller(deps, nil)
	if stopped != 11 {
		t.Fatalf("expected active poller to stop when cleared, got %d", stopped)
	}
}

func TestTelegramClearsRuntimeWhenTokenIsEmpty(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{Telegram: config.Telegram{BotToken: "old-token"}}
	currentBot := &tgbotapi.BotAPI{}
	stopped := 0

	Telegram(Deps{
		Config: cfg,
		AuthModel: phase6TelegramAuthModelStub{
			findOneByMethod: func(context.Context, string) (*modelauth.Auth, error) {
				enabled := true
				return &modelauth.Auth{
					Method:  "telegram",
					Enabled: &enabled,
					Config:  (&modelauth.TelegramAuthConfig{}).Marshal(),
				}, nil
			},
		},
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			currentBot = bot
		},
		SwapTelegramPoller: func(next func()) func() {
			prev := func() { stopped++ }
			return prev
		},
	})

	if cfg.Telegram.BotToken != "" {
		t.Fatalf("expected telegram token to clear on empty auth config, got %q", cfg.Telegram.BotToken)
	}
	if currentBot != nil {
		t.Fatal("expected telegram runtime bot to clear on empty auth config")
	}
	if stopped != 1 {
		t.Fatalf("expected existing poller to stop on empty auth config, got %d", stopped)
	}
}

type phase6TelegramAuthModelStub struct {
	findOneByMethod func(context.Context, string) (*modelauth.Auth, error)
}

func (s phase6TelegramAuthModelStub) Insert(context.Context, *modelauth.Auth) error {
	panic("unexpected Insert")
}
func (s phase6TelegramAuthModelStub) FindOne(context.Context, int64) (*modelauth.Auth, error) {
	panic("unexpected FindOne")
}
func (s phase6TelegramAuthModelStub) Update(context.Context, *modelauth.Auth) error {
	panic("unexpected Update")
}
func (s phase6TelegramAuthModelStub) Delete(context.Context, int64) error { panic("unexpected Delete") }
func (s phase6TelegramAuthModelStub) Transaction(context.Context, func(*gorm.DB) error) error {
	panic("unexpected Transaction")
}
func (s phase6TelegramAuthModelStub) GetAuthListByPage(context.Context) ([]*modelauth.Auth, error) {
	panic("unexpected GetAuthListByPage")
}
func (s phase6TelegramAuthModelStub) FindOneByMethod(ctx context.Context, method string) (*modelauth.Auth, error) {
	if s.findOneByMethod == nil {
		panic("unexpected FindOneByMethod")
	}
	return s.findOneByMethod(ctx, method)
}
func (s phase6TelegramAuthModelStub) FindAll(context.Context) ([]*modelauth.Auth, error) {
	panic("unexpected FindAll")
}
