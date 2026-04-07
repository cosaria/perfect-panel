package user

import (
	"context"
	"fmt"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"time"
)

type BindTelegramOutput struct {
	Body *types.BindTelegramResponse
}

func BindTelegramHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*BindTelegramOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*BindTelegramOutput, error) {
		l := NewBindTelegramLogic(ctx, svcCtx)
		resp, err := l.BindTelegram()
		if err != nil {
			return nil, err
		}
		return &BindTelegramOutput{Body: resp}, nil
	}
}

type BindTelegramLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Bind Telegram
func NewBindTelegramLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindTelegramLogic {
	return &BindTelegramLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindTelegramLogic) BindTelegram() (resp *types.BindTelegramResponse, err error) {
	session := l.ctx.Value("session").(string)
	return &types.BindTelegramResponse{
		Url:       fmt.Sprintf("https://t.me/%s?start=%s", l.svcCtx.Config.Telegram.BotName, session),
		ExpiredAt: time.Now().Add(300 * time.Second).UnixMilli(),
	}, nil
}
