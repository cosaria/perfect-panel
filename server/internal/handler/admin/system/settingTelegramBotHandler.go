// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
)

func SettingTelegramBotHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := system.NewSettingTelegramBotLogic(ctx, svcCtx)
		if err := l.SettingTelegramBot(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
