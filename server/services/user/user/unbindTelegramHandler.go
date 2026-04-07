// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
)

func UnbindTelegramHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewUnbindTelegramLogic(ctx, svcCtx)
		if err := l.UnbindTelegram(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
