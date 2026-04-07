package system

import (
	"context"

	"github.com/perfect-panel/server/initialize"

	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/svc"
)

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
