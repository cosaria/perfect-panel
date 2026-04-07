package system

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
)

func SettingTelegramBotHandler(deps Deps) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewSettingTelegramBotLogic(ctx, deps)
		if err := l.SettingTelegramBot(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type SettingTelegramBotLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewSettingTelegramBotLogic setting telegram bot
func NewSettingTelegramBotLogic(ctx context.Context, deps Deps) *SettingTelegramBotLogic {
	return &SettingTelegramBotLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *SettingTelegramBotLogic) SettingTelegramBot() error {
	if l.deps.ReloadTelegram != nil {
		l.deps.ReloadTelegram()
	}
	return nil
}
