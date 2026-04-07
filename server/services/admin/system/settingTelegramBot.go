package system

import (
	"context"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
)

func SettingTelegramBotHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewSettingTelegramBotLogic(ctx, svcCtx)
		if err := l.SettingTelegramBot(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type SettingTelegramBotLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSettingTelegramBotLogic setting telegram bot
func NewSettingTelegramBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettingTelegramBotLogic {
	return &SettingTelegramBotLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SettingTelegramBotLogic) SettingTelegramBot() error {
	initialize.Telegram(l.svcCtx)
	return nil
}
