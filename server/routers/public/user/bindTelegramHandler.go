// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type BindTelegramOutput struct {
	Body *types.BindTelegramResponse
}

func BindTelegramHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*BindTelegramOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*BindTelegramOutput, error) {
		l := user.NewBindTelegramLogic(ctx, svcCtx)
		resp, err := l.BindTelegram()
		if err != nil {
			return nil, err
		}
		return &BindTelegramOutput{Body: resp}, nil
	}
}
