// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateLogSettingInput struct {
	Body types.LogSetting
}

func UpdateLogSettingHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateLogSettingInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateLogSettingInput) (*struct{}, error) {
		l := log.NewUpdateLogSettingLogic(ctx, svcCtx)
		if err := l.UpdateLogSetting(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
