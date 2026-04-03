// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
)

func UnbindTelegramHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := user.NewUnbindTelegramLogic(ctx, svcCtx)
		if err := l.UnbindTelegram(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
