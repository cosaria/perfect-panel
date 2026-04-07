package user

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/types"
	"time"
)

type BindTelegramOutput struct {
	Body *types.BindTelegramResponse
}

func BindTelegramHandler(deps Deps) func(context.Context, *struct{}) (*BindTelegramOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*BindTelegramOutput, error) {
		l := NewBindTelegramLogic(ctx, deps)
		resp, err := l.BindTelegram()
		if err != nil {
			return nil, err
		}
		return &BindTelegramOutput{Body: resp}, nil
	}
}

type BindTelegramLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Bind Telegram
func NewBindTelegramLogic(ctx context.Context, deps Deps) *BindTelegramLogic {
	return &BindTelegramLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *BindTelegramLogic) BindTelegram() (resp *types.BindTelegramResponse, err error) {
	session := l.ctx.Value("session").(string)
	cfg := l.deps.currentConfig()
	return &types.BindTelegramResponse{
		Url:       fmt.Sprintf("https://t.me/%s?start=%s", cfg.Telegram.BotName, session),
		ExpiredAt: time.Now().Add(300 * time.Second).UnixMilli(),
	}, nil
}
