// huma:migrated
package system

import (
	"context"
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
